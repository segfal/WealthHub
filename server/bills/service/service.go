package service

import (
	"context"
	"server/bills/repository"
	"server/types"
	"time"
)

type Service interface {
	GetBillTotals(ctx context.Context, accountID string, startDate, endDate time.Time) (map[string]float64, error)
	GetRecurringBills(ctx context.Context, accountID string) ([]types.RecurringBill, error)
	GetUpcomingBills(ctx context.Context, accountID string) ([]types.UpcomingBill, error)
	GetBillHistory(ctx context.Context, accountID string, merchantName string) ([]types.Transaction, error)
	GetBillsByMonth(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetBillTotals(ctx context.Context, accountID string, startDate, endDate time.Time) (map[string]float64, error) {
	return s.repo.GetBillTotals(ctx, accountID, startDate, endDate)
}

func (s *service) GetRecurringBills(ctx context.Context, accountID string) ([]types.RecurringBill, error) {
	return s.repo.GetRecurringBills(ctx, accountID)
}

func (s *service) GetUpcomingBills(ctx context.Context, accountID string) ([]types.UpcomingBill, error) {
	return s.repo.GetUpcomingBills(ctx, accountID)
}

func (s *service) GetBillHistory(ctx context.Context, accountID string, merchantName string) ([]types.Transaction, error) {
	return s.repo.GetBillHistory(ctx, accountID, merchantName)
}

func (s *service) GetBillsByMonth(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error) {
	return s.repo.GetBillsByMonth(ctx, accountID, year, month)
} 