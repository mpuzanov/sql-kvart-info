package service

import (
	"context"
	"kvart-info/internal/model"
)

//go:generate mockery --name repository
type repository interface {
	Get(ctx context.Context) ([]model.SummaryInfo, error)
}

// Usecase ...
type Usecase struct {
	repo repository
}

// New ...
func New(r repository) *Usecase {
	return &Usecase{
		repo: r,
	}
}

// GetSummaryInfo получаем данные
func (uc *Usecase) GetSummaryInfo(ctx context.Context) ([]model.SummaryInfo, error) {

	data, err := uc.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}
