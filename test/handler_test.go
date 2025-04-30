package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"voucher-api/internal/handler"
	"voucher-api/internal/model"
	"voucher-api/internal/repository"
	"voucher-api/internal/service"
)

// MockTx simulates a database transaction for testing.
type MockTx struct {
	committed  bool
	rolledBack bool
}

func (m *MockTx) Commit() error {
	m.committed = true
	return nil
}

func (m *MockTx) Rollback() error {
	m.rolledBack = true
	return nil
}

func (m *MockTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return &mockResult{}, nil
}

func (m *MockTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	// Simulate a row for CreateTransaction
	return &sql.Row{}
}

// mockResult simulates a sql.Result for testing.
type mockResult struct{}

func (m *mockResult) LastInsertId() (int64, error) { return 1, nil }
func (m *mockResult) RowsAffected() (int64, error) { return 1, nil }

// MockRepository implements RepositoryInterface for testing.
type MockRepository struct{}

func (m *MockRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (repository.TxInterface, error) {
	return &MockTx{}, nil
}

func (m *MockRepository) CreateBrand(ctx context.Context, brand *model.Brand) error {
	brand.ID = 1
	brand.CreatedAt = "2025-04-30T00:00:00Z"
	return nil
}

func (m *MockRepository) CreateVoucher(ctx context.Context, voucher *model.Voucher) error {
	voucher.ID = 1
	voucher.CreatedAt = "2025-04-30T00:00:00Z"
	return nil
}

func (m *MockRepository) GetVoucher(ctx context.Context, id int) (*model.Voucher, error) {
	if id == 999 {
		return nil, sql.ErrNoRows
	}
	return &model.Voucher{
		ID:           id,
		BrandID:      1,
		Code:         "VOUCHER1",
		Name:         "Test Voucher",
		CostInPoints: 50000,
		ExpiryDate:   "2025-12-31",
		CreatedAt:    "2025-04-30T00:00:00Z",
	}, nil
}

func (m *MockRepository) GetVouchersByBrand(ctx context.Context, brandID int) ([]*model.Voucher, error) {
	if brandID == 999 {
		return []*model.Voucher{}, nil
	}
	return []*model.Voucher{{
		ID:           1,
		BrandID:      brandID,
		Code:         "VOUCHER1",
		Name:         "Test Voucher",
		CostInPoints: 50000,
		ExpiryDate:   "2025-12-31",
		CreatedAt:    "2025-04-30T00:00:00Z",
	}}, nil
}

func (m *MockRepository) CreateTransaction(ctx context.Context, tx repository.TxInterface, transaction *model.Transaction) error {
	transaction.ID = 1
	transaction.CreatedAt = "2025-04-30T00:00:00Z"
	return nil
}

func (m *MockRepository) CreateTransactionVoucher(ctx context.Context, tx repository.TxInterface, transactionID int, tv model.TransactionVoucher) error {
	return nil
}

func (m *MockRepository) GetTransaction(ctx context.Context, id int) (*model.Transaction, error) {
	if id == 999 {
		return nil, sql.ErrNoRows
	}
	return &model.Transaction{
		ID:          id,
		CustomerID:  1,
		TotalPoints: 100000,
		Vouchers:    []model.TransactionVoucher{{VoucherID: 1, Quantity: 2}},
		CreatedAt:   "2025-04-30T00:00:00Z",
	}, nil
}

func (m *MockRepository) GetCustomer(ctx context.Context, id int) (*model.Customer, error) {
	if id == 999 {
		return nil, sql.ErrNoRows
	}
	return &model.Customer{
		ID:            id,
		Name:          "John Doe",
		Email:         "john@example.com",
		PointsBalance: 200000,
		CreatedAt:     "2025-04-30T00:00:00Z",
	}, nil
}

func (m *MockRepository) UpdateCustomerPoints(ctx context.Context, tx repository.TxInterface, customerID, points int) error {
	return nil
}

func (m *MockRepository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	customer.ID = 1
	customer.CreatedAt = "2025-04-30T00:00:00Z"
	return nil
}

func TestCreateBrand(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	brand := model.Brand{Name: "Test Brand", Description: "Test Description"}
	body, _ := json.Marshal(brand)
	req, _ := http.NewRequest("POST", "/brand", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateBrand(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, status)
	}

	var response model.Brand
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != 1 || response.Name != "Test Brand" {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestCreateBrandInvalidInput(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	brand := model.Brand{Name: ""} // Invalid: empty name
	body, _ := json.Marshal(brand)
	req, _ := http.NewRequest("POST", "/brand", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateBrand(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
	}
}

func TestCreateVoucher(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	voucher := model.Voucher{
		BrandID:      1,
		Code:         "VOUCHER1",
		Name:         "Test Voucher",
		CostInPoints: 50000,
		ExpiryDate:   "2025-12-31",
	}
	body, _ := json.Marshal(voucher)
	req, _ := http.NewRequest("POST", "/voucher", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateVoucher(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, status)
	}

	var response model.Voucher
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != 1 || response.Code != "VOUCHER1" {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestCreateVoucherInvalidInput(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	voucher := model.Voucher{
		BrandID: 0, // Invalid: zero brand ID
		Code:    "",
		Name:    "",
	}
	body, _ := json.Marshal(voucher)
	req, _ := http.NewRequest("POST", "/voucher", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateVoucher(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
	}
}

func TestGetVoucher(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	req, _ := http.NewRequest("GET", "/voucher?id=1", nil)
	rr := httptest.NewRecorder()

	h.GetVoucher(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	var response model.Voucher
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != 1 || response.Code != "VOUCHER1" {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestGetVoucherNotFound(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	req, _ := http.NewRequest("GET", "/voucher?id=999", nil)
	rr := httptest.NewRecorder()

	h.GetVoucher(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
	}
}

func TestCreateCustomer(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	customer := model.Customer{
		Name:          "John Doe",
		Email:         "john@example.com",
		PointsBalance: 200000,
	}
	body, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customer", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateCustomer(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, status)
	}

	var response model.Customer
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != 1 || response.Name != "John Doe" {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestCreateCustomerInvalidInput(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	customer := model.Customer{
		Name:  "",
		Email: "", // Invalid: empty name and email
	}
	body, _ := json.Marshal(customer)
	req, _ := http.NewRequest("POST", "/customer", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.CreateCustomer(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
	}
}

func TestGetVouchersByBrand(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	req, _ := http.NewRequest("GET", "/voucher/brand?id=1", nil)
	rr := httptest.NewRecorder()

	h.GetVouchersByBrand(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	var response []*model.Voucher
	json.NewDecoder(rr.Body).Decode(&response)
	if len(response) != 1 || response[0].ID != 1 || response[0].BrandID != 1 {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestGetVouchersByBrandEmpty(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	req, _ := http.NewRequest("GET", "/voucher/brand?id=999", nil)
	rr := httptest.NewRecorder()

	h.GetVouchersByBrand(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	var response []*model.Voucher
	json.NewDecoder(rr.Body).Decode(&response)
	if len(response) != 0 {
		t.Errorf("Expected empty response, got %+v", response)
	}
}

func TestMakeRedemption(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	redemption := model.RedemptionRequest{
		CustomerID: 1,
		Vouchers:   []model.TransactionVoucher{{VoucherID: 1, Quantity: 2}},
	}
	body, _ := json.Marshal(redemption)
	req, _ := http.NewRequest("POST", "/transaction/redemption", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.MakeRedemption(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %v, got %v", http.StatusCreated, status)
	}

	var response model.Transaction
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != 1 || response.CustomerID != 1 || response.TotalPoints != 100000 {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestMakeRedemptionInvalidCustomer(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	redemption := model.RedemptionRequest{
		CustomerID: 999, // Non-existent customer
		Vouchers:   []model.TransactionVoucher{{VoucherID: 1, Quantity: 2}},
	}
	body, _ := json.Marshal(redemption)
	req, _ := http.NewRequest("POST", "/transaction/redemption", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.MakeRedemption(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
	}
}

func TestMakeRedemptionInvalidVoucher(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	redemption := model.RedemptionRequest{
		CustomerID: 1,
		Vouchers:   []model.TransactionVoucher{{VoucherID: 999, Quantity: 2}}, // Non-existent voucher
	}
	body, _ := json.Marshal(redemption)
	req, _ := http.NewRequest("POST", "/transaction/redemption", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.MakeRedemption(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
	}
}

func TestMakeRedemptionInsufficientPoints(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	redemption := model.RedemptionRequest{
		CustomerID: 1,
		Vouchers:   []model.TransactionVoucher{{VoucherID: 1, Quantity: 10}}, // Requires 500,000 points
	}
	body, _ := json.Marshal(redemption)
	req, _ := http.NewRequest("POST", "/transaction/redemption", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	h.MakeRedemption(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %v, got %v", http.StatusBadRequest, status)
	}
}

func TestGetTransactionDetail(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	req, _ := http.NewRequest("GET", "/transaction/redemption?transactionId=1", nil)
	rr := httptest.NewRecorder()

	h.GetTransactionDetail(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	var response model.Transaction
	json.NewDecoder(rr.Body).Decode(&response)
	if response.ID != 1 || response.CustomerID != 1 {
		t.Errorf("Unexpected response: %+v", response)
	}
}

func TestGetTransactionDetailNotFound(t *testing.T) {
	repo := &MockRepository{}
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	req, _ := http.NewRequest("GET", "/transaction/redemption?transactionId=999", nil)
	rr := httptest.NewRecorder()

	h.GetTransactionDetail(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
	}
}
