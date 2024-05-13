package service

import (
	"context"
	"kvart-info/internal/model"
)

//go:generate mockery --name repository
type repository interface {
	GetByTip(ctx context.Context, tipID any) ([]model.SummaryInfo, error)
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
	var tipID any = nil
	data, err := uc.repo.GetByTip(ctx, tipID)
	if err != nil {
		return nil, err
	}
	return data, nil
}
