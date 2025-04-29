package model

type Brand struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type Voucher struct {
	ID           int    `json:"id"`
	BrandID      int    `json:"brand_id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	CostInPoints int    `json:"cost_in_points"`
	ExpiryDate   string `json:"expiry_date"`
	CreatedAt    string `json:"created_at"`
}

type Customer struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	PointsBalance int    `json:"points_balance"`
	CreatedAt     string `json:"created_at"`
}

type Transaction struct {
	ID          int                  `json:"id"`
	CustomerID  int                  `json:"customer_id"`
	TotalPoints int                  `json:"total_points"`
	Vouchers    []TransactionVoucher `json:"vouchers"`
	CreatedAt   string               `json:"created_at"`
}

type TransactionVoucher struct {
	VoucherID int `json:"voucher_id"`
	Quantity  int `json:"quantity"`
}

type RedemptionRequest struct {
	CustomerID int                  `json:"customer_id"`
	Vouchers   []TransactionVoucher `json:"vouchers"`
}
