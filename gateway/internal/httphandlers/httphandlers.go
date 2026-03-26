package httphandlers

import (
	"encoding/json"
	"gateway/internal/service"
	"net/http"
	"strconv"
)

func ProductCatalogHandler(s *service.GatewayService) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		pageNum := r.URL.Query().Get("page_number")
		paramPage, err := strconv.ParseUint(pageNum, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		itemsPerPage := r.URL.Query().Get("items_per_page")
		paramItems, err := strconv.ParseUint(itemsPerPage, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		products, err := s.ProductCatalog(r.Context(), paramPage, paramItems)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.Marshal(&products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

	return http.HandlerFunc(h)
}

func CreateUserHandler(s *service.GatewayService) http.HandlerFunc {
	type Request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := s.NewUser(r.Context(), req.Login, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"status":  "ok",
			"user_id": userId,
		}
		data, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

	return http.HandlerFunc(h)
}

func DeleteUserHandler(s *service.GatewayService) http.HandlerFunc {
	type Request struct {
		UserId uint64 `json:"user_id"`
	}

	req := Request{}
	h := func(w http.ResponseWriter, r *http.Request) {
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = s.DeleteUser(r.Context(), req.UserId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]string{"status": "ok"}
		data, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

	return http.HandlerFunc(h)
}

func BasketInfoHandler(s *service.GatewayService) http.HandlerFunc {
	h := func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")
		id, err := strconv.ParseUint(userId, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		productsInBasket, err := s.BasketInfo(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp := map[string]interface{}{
			"status": "ok",
			"data":   productsInBasket,
		}
		data, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}

	return http.HandlerFunc(h)
}

func AddProductToBasketHandler(s *service.GatewayService) http.HandlerFunc {
	type Request struct {
		UserId          uint64 `json:"user_id"`
		ProductId       uint64 `json:"product_id"`
		ProductQuantity uint64 `json:"product_quantity"`
	}

	h := func(w http.ResponseWriter, r *http.Request) {
		req := Request{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = s.AddProductToBasket(r.Context(), req.UserId, req.ProductId, req.ProductQuantity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp := map[string]string{
			"status": "ok",
		}
		data, err := json.Marshal(&resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}

	return http.HandlerFunc(h)
}
