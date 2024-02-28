package controller

import (
	"context"
	"kvart-info/internal/config"
	"kvart-info/internal/model"
	"kvart-info/pkg/logging"
)

// UseCase ...
type UseCase interface {
	GetSummaryInfo(ctx context.Context) ([]model.SummaryInfo, error)
}

// Controller ...
type Controller struct {
	cfg *config.Config
	log *logging.Logger
	uc  UseCase
}

// New ...
func New(cfg *config.Config, log *logging.Logger, uc UseCase) *Controller {
	return &Controller{cfg: cfg, log: log, uc: uc}
}
