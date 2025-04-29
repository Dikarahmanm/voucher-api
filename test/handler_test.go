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
	// Simulate a row that returns ID=1 and CreatedAt for CreateTransaction
	row := &sql.Row{}
	return row
}

// mockResult simulates a sql.Result for testing.
type mockResult struct{}

func (m *mockResult) LastInsertId() (int64, error) { return 1, nil }
func (m *mockResult) RowsAffected() (int64, error) { return 1, nil }

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
	return &model.Transaction{
		ID:          id,
		CustomerID:  1,
		TotalPoints: 100000,
		Vouchers:    []model.TransactionVoucher{{VoucherID: 1, Quantity: 2}},
		CreatedAt:   "2025-04-30T00:00:00Z",
	}, nil
}

func (m *MockRepository) GetCustomer(ctx context.Context, id int) (*model.Customer, error) {
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
