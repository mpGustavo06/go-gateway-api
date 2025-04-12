package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/domain"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/dto"
	"github.com/mpGustavo06/go-gateway-api/go-gateway/internal/service"
)

type InvoiceHandler struct {
	service *service.InvoiceService
}

func NewInvoiceHandler(service *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateInvoiceInput

	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {		
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.APIKey = r.Header.Get("X-API-Key")

	output, err := h.service.Create(input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

 func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	apiKey := r.Header.Get("X-API-Key")

	if apiKey == "" {		
		http.Error(w, "API Key header is required", http.StatusBadRequest)
		return
	}

	output, err := h.service.GetById(id, apiKey)

	if err != nil {
		switch err {
			case domain.ErrInvoiceNotFound:
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			case domain.ErrAccountNotFound:
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			case domain.ErrUnauthorizedAccess:
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}

func (h *InvoiceHandler) ListByAccount(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("X-API-Key")

	if apiKey == "" {
		http.Error(w, "API Key header is required", http.StatusUnauthorized)
		return
	}

	output, err := h.service.ListByAccountAPIKey(apiKey)

	if err != nil {
		switch err {
			case domain.ErrAccountNotFound:
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
		}																	
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}