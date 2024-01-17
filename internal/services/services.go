package services

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"kvart-info/internal/config"
	"kvart-info/internal/email"
	"kvart-info/internal/model"
	"kvart-info/pkg/logging"
	"strconv"
	"text/tabwriter"
	"time"
)

// ServiceInfo ...
type ServiceInfo struct {
	ctx context.Context
	cfg *config.Config
	db  Datastore
}

// Datastore ...
type Datastore interface {
	// GetTotalData получение сводной информации
	GetTotalData() ([]model.TotalData, error)
}

// New ...
func New(ctx context.Context, cfg *config.Config, db Datastore) *ServiceInfo {
	return &ServiceInfo{ctx: ctx, cfg: cfg, db: db}
}

// Run выполняем сервис получения данных по БД
func (s *ServiceInfo) Run() error {
	l := logging.LoggerFromContext(s.ctx)

	data, err := s.GetLicTotal()
	if err != nil {
		return err
	}

	bodyMessage, title, err := s.CreateBodyText(data)
	if err != nil {
		return err
	}

	if s.cfg.Mail.IsSendEmail {

		emailStatus := make(chan string)
		go email.New(s.cfg).Send(bodyMessage, title, s.cfg.Mail.ToSendEmail, "", emailStatus)
		status := <-emailStatus
		l.Info(status)

	} else {
		fmt.Println(bodyMessage)
	}
	return nil
}

// GetLicTotal получаем данные
func (s *ServiceInfo) GetLicTotal() ([]model.TotalData, error) {

	data, err := s.db.GetTotalData()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateBodyText формируем письмо
func (s *ServiceInfo) CreateBodyText(data []model.TotalData) (string, string, error) {

	type ViewData struct {
		Title     string
		CreatedOn string
		Body      []model.TotalData
	}
	title := fmt.Sprintf("Information about the %s database, create: %s",
		s.cfg.DB.Database,
		time.Now().Format("02.01.2006"),
	)
	var body bytes.Buffer

	if s.cfg.Mail.IsSendEmail {
		t, err := template.New("").Funcs(template.FuncMap{
			"IncInt": func(i int) string {
				i++
				return strconv.Itoa(i)
			},
		}).Parse(bodyTemplateMap)
		if err != nil {
			return "", "", err
		}

		p := &ViewData{Title: title, CreatedOn: time.Now().Format("02.01.2006"), Body: data}

		if err := t.Execute(&body, p); err != nil {
			return "", "", err
		}
		bodyMessage := body.String()

		return bodyMessage, title, nil
	}

	// текстовый шаблон таблицы
	t := template.Must(template.New("").Parse(tmplText))
	w := tabwriter.NewWriter(&body, 5, 0, 3, ' ', 0)
	if err := t.Execute(w, data); err != nil {
		return "", "", err
	}
	_ = w.Flush()
	bodyMessage := body.String()

	return bodyMessage, title, nil
}
