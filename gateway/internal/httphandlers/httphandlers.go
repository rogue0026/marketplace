package httphandlers

import (
	"encoding/json"
	"net/http"
)

func ProductCatalogHandler() http.Handler {
	type Request struct {
		PageNumber   int `json:"page_number"`
		ItemsPerPage int `json:"items_per_page"`
	}

	req := Request{}
	h := func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	return http.HandlerFunc(h)
}
