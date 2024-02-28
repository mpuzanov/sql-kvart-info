package app

import (
	"kvart-info/internal/config"
	"kvart-info/internal/controller"
	"kvart-info/internal/repository"
	"kvart-info/internal/service"
	"kvart-info/pkg/logging"
	"kvart-info/pkg/mssql"

	"github.com/pkg/errors"
)

// Run ...
func Run(cfg *config.Config) error {

	logger := logging.NewLogger(cfg.Env)

	logger.Debug("debug", "cfg", cfg)

	mssql, err := mssql.New(&cfg.DB)
	if err != nil {
		return errors.Wrap(err, "mssql.New")
	}
	defer mssql.Close()

	repoInfo := repository.New(mssql)

	infoUseCase := service.New(repoInfo)

	c := controller.New(cfg, logger, infoUseCase)

	return c.OutputInfo()
}
