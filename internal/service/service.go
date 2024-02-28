package service

import (
	"context"
	"kvart-info/internal/model"
)

// SummaryInfoService ...
type SummaryInfoService interface {
	GetSummaryInfo(ctx context.Context) ([]model.SummaryInfo, error)
}

// repositorySummaryInfo ...
type repository interface {
	Get(ctx context.Context) ([]model.SummaryInfo, error)
}

// UseCase ...
type UseCase struct {
	repo repository
}

// New ...
func New(r repository) *UseCase {
	return &UseCase{
		repo: r,
	}
}

// GetSummaryInfo получаем данные
func (uc *UseCase) GetSummaryInfo(ctx context.Context) ([]model.SummaryInfo, error) {

	data, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
