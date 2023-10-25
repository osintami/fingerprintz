// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"

	"github.com/osintami/fingerprintz/common"
)

func (x *NormalizedDataServer) GetItemHandler(w http.ResponseWriter, r *http.Request) {
	x.PostItemHandler(w, r)
}
func (x *NormalizedDataServer) PostItemHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := x.params.Keys(r)
	if len(keys) == 0 || err != nil {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}
	x.ItemHandler(w, r, keys)
}
func (x *NormalizedDataServer) ItemHandler(w http.ResponseWriter, r *http.Request, keys map[string]string) {
	categoryName := common.PathParam(r, "category")
	vendorName := common.PathParam(r, "vendor")
	itemName := common.PathParam(r, "item")

	// itemName, err := url.QueryUnescape(itemName)
	// if err != nil {
	// 	common.SendError(w, ErrInvalidItemParam, http.StatusBadRequest)
	// 	return
	// }

	if categoryName == "" {
		common.SendError(w, ErrInvalidCategoryParam, http.StatusBadRequest)
		return
	}
	if vendorName == "" {
		common.SendError(w, ErrInvalidSourceParam, http.StatusBadRequest)
		return
	}
	if itemName == "" {
		common.SendError(w, ErrInvalidItemParam, http.StatusBadRequest)
		return
	}

	if categoryName != CATEGORY_RULE && keys[categoryName] == "" {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}

	// build item
	item := NewDataURI(categoryName, vendorName, itemName)

	// used by some of the fancier sources for internal routing
	keys[common.INPUT_RULE] = item.Key()
	//	keys[common.INPUT_TYPE] = "Boolean"

	var output *common.DataOutput
	output, err := x.router.DataValue(r.Context(), item, keys)

	// for some errors, return data
	if err != nil && err != common.ErrNoDataPresent {
		common.SendError(w, err, http.StatusNotFound)
		return
	}
	// strip out internal use only keys from response
	delete(output.Keys, common.INPUT_ROLE)
	delete(output.Keys, common.INPUT_KEY)
	delete(output.Keys, common.INPUT_TYPE)
	delete(output.Keys, common.INPUT_RULE)

	common.SendPrettyJSON(w, output)
}
