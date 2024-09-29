package app

import (
	"fmt"
	"kvart-info/internal/config"
	"kvart-info/internal/controller"
	"kvart-info/internal/controller/notify"
	"kvart-info/internal/repository"
	"kvart-info/internal/service"
	"os"

	"github.com/mpuzanov/wslog"

	"github.com/mpuzanov/dbwrap"
)

// Run ...
func Run(cfg *config.Config) error {

	logger := wslog.NewEnv(cfg.Env)
	logger.Debug("debug", wslog.Any("cfg", cfg))

	dbCon, err := dbwrap.NewConnect(&cfg.DB)
	if err != nil {
		return fmt.Errorf("dbwrap.New : %w", err)
	}
	defer dbCon.Close()
	logger.Debug("DB", "cfg.DB", cfg.DB.String())

	repoInfo := repository.New(dbCon)

	infoUseCase := service.New(repoInfo)

	bodyMessage, title, err := controller.New(cfg, logger, infoUseCase).OutputInfo()
	if err != nil {
		return fmt.Errorf("controller OutputInfo : %w", err)
	}

	if cfg.IsSendEmail && cfg.ToSendEmail != "" {

		objNotify := notify.NewEmail(cfg.Mail, bodyMessage, title, cfg.ToSendEmail)
		emailStatus, err := objNotify.Send()
		if err != nil {
			return fmt.Errorf("NotifyEmail : %w", err)
		}
		logger.Info(emailStatus)
		return nil
	}

	// выдаём на экран
	notify.NewCli(bodyMessage, os.Stdout).Send()

	return nil
}
