package repository

import (
	"context"
	"database/sql"
	"voucher-api/internal/model"
)

// TxInterface defines the methods for a transaction.
type TxInterface interface {
	Commit() error
	Rollback() error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// RepositoryInterface defines the methods for interacting with the database.
type RepositoryInterface interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (TxInterface, error)
	CreateBrand(ctx context.Context, brand *model.Brand) error
	CreateVoucher(ctx context.Context, voucher *model.Voucher) error
	GetVoucher(ctx context.Context, id int) (*model.Voucher, error)
	GetVouchersByBrand(ctx context.Context, brandID int) ([]*model.Voucher, error)
	CreateTransaction(ctx context.Context, tx TxInterface, transaction *model.Transaction) error
	CreateTransactionVoucher(ctx context.Context, tx TxInterface, transactionID int, tv model.TransactionVoucher) error
	GetTransaction(ctx context.Context, id int) (*model.Transaction, error)
	GetCustomer(ctx context.Context, id int) (*model.Customer, error)
	UpdateCustomerPoints(ctx context.Context, tx TxInterface, customerID, points int) error
	CreateCustomer(ctx context.Context, customer *model.Customer) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) BeginTx(ctx context.Context, opts *sql.TxOptions) (TxInterface, error) {
	tx, err := r.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *Repository) CreateBrand(ctx context.Context, brand *model.Brand) error {
	query := `INSERT INTO brands (name, description) VALUES ($1, $2) RETURNING id, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`
	return r.db.QueryRowContext(ctx, query, brand.Name, brand.Description).Scan(&brand.ID, &brand.CreatedAt)
}

func (r *Repository) CreateVoucher(ctx context.Context, voucher *model.Voucher) error {
	query := `INSERT INTO vouchers (brand_id, code, name, cost_in_points, expiry_date) 
             VALUES ($1, $2, $3, $4, $5) RETURNING id, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`
	return r.db.QueryRowContext(ctx, query, voucher.BrandID, voucher.Code, voucher.Name,
		voucher.CostInPoints, voucher.ExpiryDate).Scan(&voucher.ID, &voucher.CreatedAt)
}

func (r *Repository) GetVoucher(ctx context.Context, id int) (*model.Voucher, error) {
	voucher := &model.Voucher{}
	query := `SELECT id, brand_id, code, name, cost_in_points, expiry_date, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') 
             FROM vouchers WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&voucher.ID, &voucher.BrandID, &voucher.Code, &voucher.Name,
		&voucher.CostInPoints, &voucher.ExpiryDate, &voucher.CreatedAt)
	if err != nil {
		return nil, err
	}
	return voucher, nil
}

func (r *Repository) GetVouchersByBrand(ctx context.Context, brandID int) ([]*model.Voucher, error) {
	query := `SELECT id, brand_id, code, name, cost_in_points, expiry_date, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') 
             FROM vouchers WHERE brand_id = $1`
	rows, err := r.db.QueryContext(ctx, query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vouchers []*model.Voucher
	for rows.Next() {
		voucher := &model.Voucher{}
		if err := rows.Scan(&voucher.ID, &voucher.BrandID, &voucher.Code, &voucher.Name,
			&voucher.CostInPoints, &voucher.ExpiryDate, &voucher.CreatedAt); err != nil {
			return nil, err
		}
		vouchers = append(vouchers, voucher)
	}
	return vouchers, nil
}

func (r *Repository) CreateTransaction(ctx context.Context, tx TxInterface, transaction *model.Transaction) error {
	query := `INSERT INTO transactions (customer_id, total_points) VALUES ($1, $2) 
             RETURNING id, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`
	return tx.QueryRowContext(ctx, query, transaction.CustomerID, transaction.TotalPoints).
		Scan(&transaction.ID, &transaction.CreatedAt)
}

func (r *Repository) CreateTransactionVoucher(ctx context.Context, tx TxInterface, transactionID int, tv model.TransactionVoucher) error {
	query := `INSERT INTO transaction_vouchers (transaction_id, voucher_id, quantity) 
             VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, query, transactionID, tv.VoucherID, tv.Quantity)
	return err
}

func (r *Repository) GetTransaction(ctx context.Context, id int) (*model.Transaction, error) {
	transaction := &model.Transaction{}
	query := `SELECT id, customer_id, total_points, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') 
             FROM transactions WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID, &transaction.CustomerID, &transaction.TotalPoints, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}

	query = `SELECT voucher_id, quantity FROM transaction_vouchers WHERE transaction_id = $1`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tv model.TransactionVoucher
		if err := rows.Scan(&tv.VoucherID, &tv.Quantity); err != nil {
			return nil, err
		}
		transaction.Vouchers = append(transaction.Vouchers, tv)
	}
	return transaction, nil
}

func (r *Repository) GetCustomer(ctx context.Context, id int) (*model.Customer, error) {
	customer := &model.Customer{}
	query := `SELECT id, name, email, points_balance, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') 
             FROM customers WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID, &customer.Name, &customer.Email, &customer.PointsBalance, &customer.CreatedAt)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *Repository) UpdateCustomerPoints(ctx context.Context, tx TxInterface, customerID, points int) error {
	query := `UPDATE customers SET points_balance = points_balance - $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, points, customerID)
	return err
}

func (r *Repository) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	query := `INSERT INTO customers (name, email, points_balance) 
             VALUES ($1, $2, $3) RETURNING id, to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')`
	return r.db.QueryRowContext(ctx, query, customer.Name, customer.Email, customer.PointsBalance).
		Scan(&customer.ID, &customer.CreatedAt)
}
