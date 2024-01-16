package services

import (
	"bytes"
	"fmt"
	"html/template"
	"kvart-info/internal/config"
	"kvart-info/internal/email"
	"kvart-info/internal/model"
	"kvart-info/internal/storage"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// Info ...
type Info struct {
	cfg *config.Config
	db  storage.Datastore
}

// InfoData ...
type InfoData interface {
	GetLicTotal(isSendMail bool) error
}

func New(cfg *config.Config, db storage.Datastore) *Info {
	return &Info{cfg: cfg, db: db}
}

// Run выполняем сервис получения данных по БД
func (s *Info) Run() error {

	data, err := s.GetLicTotal()
	if err != nil {
		return err
	}

	bodyMessage, title, err := s.CreateBodyText(data)
	if err != nil {
		return err
	}

	if s.cfg.Mail.IsSendEmail {
		err = email.New(s.cfg).Send(bodyMessage, title, s.cfg.Mail.ToSendEmail, "")
		if err != nil {
			return err
		}
	} else {
		fmt.Println(bodyMessage)
	}
	return nil
}

// GetLicTotal получаем данные
func (s *Info) GetLicTotal() ([]model.TotalData, error) {

	slog.Info("Executing query", "database", s.cfg.DB.Database)

	data, err := s.db.GetTotalData()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateBodyText формируем письмо
func (s *Info) CreateBodyText(data []model.TotalData) (string, string, error) {

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

	if strings.ToLower(s.cfg.Mail.ContentType) == "text/html" {
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
