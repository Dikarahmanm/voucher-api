package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

func (h *Handler) GetTransactionDetail(w http.ResponseWriter, r *http.Request) {
	transactionIDStr := r.URL.Query().Get("transactionId")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil || transactionID <= 0 {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.svc.GetTransactionDetail(r.Context(), transactionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}
