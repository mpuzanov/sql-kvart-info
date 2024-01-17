package services

import (
	"bytes"
	"fmt"
	"html/template"
	"kvart-info/internal/config"
	"kvart-info/internal/email"
	"kvart-info/internal/model"
	"log"
	"log/slog"
	"strconv"
	"text/tabwriter"
	"time"
)

// serviceInfo ...
type serviceInfo struct {
	cfg *config.Config
	db  Datastore
}

type Datastore interface {
	// GetTotalData получение сводной информации
	GetTotalData() ([]model.TotalData, error)
}

func New(cfg *config.Config, db Datastore) *serviceInfo {
	return &serviceInfo{cfg: cfg, db: db}
}

// Run выполняем сервис получения данных по БД
func (s *serviceInfo) Run() error {

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
		slog.Info(status)

	} else {
		fmt.Println(bodyMessage)
	}
	return nil
}

// GetLicTotal получаем данные
func (s *serviceInfo) GetLicTotal() ([]model.TotalData, error) {

	slog.Info("Executing query", "database", s.cfg.DB.Database)

	data, err := s.db.GetTotalData()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateBodyText формируем письмо
func (s *serviceInfo) CreateBodyText(data []model.TotalData) (string, string, error) {

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
				i += 1
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
	itemsTmpl := template.Must(template.New("").Parse(tmplText))
	w := tabwriter.NewWriter(&body, 5, 0, 3, ' ', 0)
	if err := itemsTmpl.Execute(w, data); err != nil {
		log.Fatal("error executing items template")
	}
	_ = w.Flush()
	bodyMessage := body.String()

	return bodyMessage, title, nil
}
