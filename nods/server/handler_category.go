// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"fmt"
	"net/http"
	"strconv"
	"text/tabwriter"

	"github.com/osintami/fingerprintz/common"
)

func (x *NormalizedDataServer) GetCategoryHandler(w http.ResponseWriter, r *http.Request) {
	x.PostCategoryHandler(w, r)
}
func (x *NormalizedDataServer) PostCategoryHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := x.params.Keys(r)
	if len(keys) == 0 || err != nil {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}
	x.CategoryHandler(w, r, keys)
}
func (x *NormalizedDataServer) CategoryHandler(w http.ResponseWriter, r *http.Request, keys map[string]string) {
	categoryName := common.PathParam(r, "category")
	if categoryName == CATEGORY_RULE || categoryName == "" {
		common.SendError(w, ErrInvalidCategoryParam, http.StatusBadRequest)
		return
	}

	results, err := x.router.CategoryValues(r.Context(), categoryName, keys)
	if err != nil {
		common.SendError(w, err, http.StatusBadRequest)
		return
	}

	wantCSV, _ := strconv.ParseBool(keys["csv"])
	if !wantCSV {
		common.SendJSON(w, results)
	} else {
		w := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.TabIndent)
		_, _ = fmt.Fprintf(w, "category,\tsource,\titem,\tkey,\tvalue\n")

		for _, result := range results {
			if result != nil {
				uri := NewItemSplitter(result.Item)
				_, _ = fmt.Fprintf(w, "%s,\t%s,\t%s,\t%s,\t%s\n",
					uri.CategoryName,
					uri.SourceName,
					uri.ItemName,
					keys[categoryName],
					result.Result.Raw)
			}
		}
		_ = w.Flush()
	}
}
