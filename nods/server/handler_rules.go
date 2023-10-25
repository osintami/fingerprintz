// Copyright © 2023 OSINTAMI. This is not yours.
package server

import (
	"net/http"

	"github.com/osintami/fingerprintz/common"
)

func (x *NormalizedDataServer) GetEvaluateHandler(w http.ResponseWriter, r *http.Request) {
	x.PostEvaluateHandler(w, r)
}
func (x *NormalizedDataServer) PostEvaluateHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := x.params.Keys(r)
	if len(keys) == 0 || err != nil {
		common.SendError(w, ErrMissingInputs, http.StatusBadRequest)
		return
	}
	keys[common.INPUT_TYPE] = "JSON"
	x.EvaluateHandler(w, r, keys)
}

func (x NormalizedDataServer) EvaluateHandler(w http.ResponseWriter, r *http.Request, keys common.DataInputs) {
	value, err := x.rules.CategoryInfo(r.Context(), CATEGORY_RULE, keys)
	if err != nil {
		common.SendError(w, err, http.StatusInternalServerError)
		return
	}
	common.SendJSON(w, value)
}
