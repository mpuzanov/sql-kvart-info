package app

import (
	"kvart-info/internal/config"
	"kvart-info/internal/controller"
	"kvart-info/internal/controller/notify"
	"kvart-info/internal/repository"
	"kvart-info/internal/service"
	"kvart-info/pkg/logging"
	"kvart-info/pkg/mssql"
	"os"

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

	bodyMessage, title, err := controller.New(cfg, logger, infoUseCase).OutputInfo()
	if err != nil {
		return errors.Wrap(err, "OutputInfo")
	}

	if cfg.IsSendEmail && cfg.ToSendEmail != "" {

		objNotify := notify.NotifyEmail{
			Cfg:         cfg.Mail,
			BodyMessage: bodyMessage,
			Title:       title,
			ToSendEmail: cfg.ToSendEmail,
		}
		emailStatus, err := notify.New(objNotify).Send()
		if err != nil {
			return errors.Wrap(err, "NotifyEmail")
		}
		logger.Info(emailStatus)
		return nil
	}

	// выдаём на экран
	objNotify := notify.NotifyCli{BodyMessage: bodyMessage, Writer: os.Stdout}
	notify.New(objNotify).Send()

	return nil
}
