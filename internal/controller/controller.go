package controller

import (
	"context"
	"kvart-info/internal/config"
	"kvart-info/internal/model"
	"kvart-info/pkg/logging"
)

// useсase ...
//
//go:generate mockery --name useсase
type usecase interface {
	GetSummaryInfo(ctx context.Context) ([]model.SummaryInfo, error)
}

// Controller ...
type Controller struct {
	cfg *config.Config
	log *logging.Logger
	uc  usecase
}

// New ...
func New(cfg *config.Config, log *logging.Logger, uc usecase) *Controller {
	return &Controller{cfg: cfg, log: log, uc: uc}
}