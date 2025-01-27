package service

import (
	"context"
	"server/income/repository"
	"server/types"
)

type Service interface {
	GetIncome(ctx context.Context, accountID string) ([]types.Transaction, error)
	GetMonthlyIncome(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetIncome(ctx context.Context, accountID string) ([]types.Transaction, error) {
	return s.repo.GetIncome(ctx, accountID)
}

func (s *service) GetMonthlyIncome(ctx context.Context, accountID string, year int, month int) ([]types.Transaction, error) {
	return s.repo.GetMonthlyIncome(ctx, accountID, year, month)
} 