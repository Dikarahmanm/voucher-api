package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"voucher-api/internal/model"
	"voucher-api/internal/service"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var brand model.Brand
	if err := json.NewDecoder(r.Body).Decode(&brand); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateBrand(r.Context(), &brand); err != nil {
		if strings.Contains(err.Error(), "brand name is required") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create brand: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(brand)
}

func (h *Handler) CreateVoucher(w http.ResponseWriter, r *http.Request) {
	var voucher model.Voucher
	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateVoucher(r.Context(), &voucher); err != nil {
		if strings.Contains(err.Error(), "invalid voucher data") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create voucher: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(voucher)
}

func (h *Handler) GetVoucher(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid voucher ID", http.StatusBadRequest)
		return
	}

	voucher, err := h.svc.GetVoucher(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Voucher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get voucher: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(voucher)
}

func (h *Handler) GetVouchersByBrand(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid brand ID", http.StatusBadRequest)
		return
	}

	vouchers, err := h.svc.GetVouchersByBrand(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to get vouchers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(vouchers)
}

func (h *Handler) MakeRedemption(w http.ResponseWriter, r *http.Request) {
	var req model.RedemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.svc.MakeRedemption(r.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid redemption request") || strings.Contains(err.Error(), "insufficient points") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if strings.Contains(err.Error(), "customer not found") || strings.Contains(err.Error(), "voucher") && strings.Contains(err.Error(), "not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to process redemption: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) GetTransactionDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("transactionId")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.svc.GetTransactionDetail(r.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Transaction not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer model.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateCustomer(r.Context(), &customer); err != nil {
		if strings.Contains(err.Error(), "invalid customer data") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create customer: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}
