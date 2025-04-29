package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"voucher-api/internal/model"
	"voucher-api/internal/repository"
)

type Service struct {
	repo repository.RepositoryInterface
}

func NewService(repo repository.RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBrand(ctx context.Context, brand *model.Brand) error {
	if brand.Name == "" {
		return errors.New("brand name is required")
	}
	return s.repo.CreateBrand(ctx, brand)
}

func (s *Service) CreateVoucher(ctx context.Context, voucher *model.Voucher) error {
	if voucher.Code == "" || voucher.Name == "" || voucher.BrandID == 0 {
		return errors.New("invalid voucher data")
	}
	return s.repo.CreateVoucher(ctx, voucher)
}

func (s *Service) GetVoucher(ctx context.Context, id int) (*model.Voucher, error) {
	return s.repo.GetVoucher(ctx, id)
}

func (s *Service) GetVouchersByBrand(ctx context.Context, brandID int) ([]*model.Voucher, error) {
	return s.repo.GetVouchersByBrand(ctx, brandID)
}

func (s *Service) MakeRedemption(ctx context.Context, req *model.RedemptionRequest) (*model.Transaction, error) {
	if req.CustomerID == 0 || len(req.Vouchers) == 0 {
		return nil, errors.New("invalid redemption request")
	}

	customer, err := s.repo.GetCustomer(ctx, req.CustomerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	var totalPoints int
	for _, tv := range req.Vouchers {
		voucher, err := s.repo.GetVoucher(ctx, tv.VoucherID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("voucher %d not found", tv.VoucherID)
			}
			return nil, err
		}
		totalPoints += voucher.CostInPoints * tv.Quantity
	}

	if customer.PointsBalance < totalPoints {
		return nil, errors.New("insufficient points")
	}

	tx, err := s.repo.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	transaction := &model.Transaction{
		CustomerID:  req.CustomerID,
		TotalPoints: totalPoints,
		Vouchers:    req.Vouchers,
	}

	if err := s.repo.CreateTransaction(ctx, tx, transaction); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, tv := range req.Vouchers {
		if err := s.repo.CreateTransactionVoucher(ctx, tx, transaction.ID, tv); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := s.repo.UpdateCustomerPoints(ctx, tx, req.CustomerID, totalPoints); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (s *Service) GetTransactionDetail(ctx context.Context, id int) (*model.Transaction, error) {
	return s.repo.GetTransaction(ctx, id)
}

func (s *Service) CreateCustomer(ctx context.Context, customer *model.Customer) error {
	if customer.Name == "" || customer.Email == "" {
		return errors.New("invalid customer data")
	}
	return s.repo.CreateCustomer(ctx, customer)
}
