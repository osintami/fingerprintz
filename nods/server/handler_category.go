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

	outputs, err := x.router.CategoryValues(r.Context(), categoryName, keys)
	if err != nil {
		common.SendError(w, err, http.StatusBadRequest)
		return
	}

	wantCSV, _ := strconv.ParseBool(keys["csv"])
	if !wantCSV {
		common.SendPrettyJSON(w, outputs)
		return
	} else {
		w := tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.TabIndent)
		for _, output := range outputs {
			if output != nil {
				value := output.Result.Raw
				if value == "" {
					value = "\"\""
				}

				fmt.Fprintf(w, "%s,\t%s\n",
					output.Item,
					value)
			}
		}
		_ = w.Flush()
	}
}
