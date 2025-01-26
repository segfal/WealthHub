package service

import (
	"context"
	"server/categories/repository"
	"server/types"
)

type Service interface {
	GetCategories(ctx context.Context, accountID string) ([]types.Category, error)
	GetCategoryTotals(ctx context.Context, accountID string) (map[string]float64, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCategories(ctx context.Context, accountID string) ([]types.Category, error) {
	return s.repo.GetCategories(ctx, accountID)
}

func (s *service) GetCategoryTotals(ctx context.Context, accountID string) (map[string]float64, error) {
	return s.repo.GetCategoryTotals(ctx, accountID)
} 