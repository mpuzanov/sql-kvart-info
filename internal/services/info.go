package services

import (
	"bytes"
	"fmt"
	"html/template"
	"kvart-info/internal/config"
	"kvart-info/internal/email"
	"kvart-info/internal/storage"
	"log/slog"
	"strconv"
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

// GetLicTotal ...
func (s *Info) GetLicTotal() error {

	slog.Info("Executing query", "database", s.cfg.DB.Database)

	data, err := s.db.Query(storage.QueryGetTotal, map[string]interface{}{})
	if err != nil {
		return err
	}

	type ViewData struct {
		Title     string
		CreatedOn string
		Body      []map[string]interface{}
	}

	t, err := template.New("").Funcs(template.FuncMap{
		"IncInt": func(i int) string {
			i += 1
			return strconv.Itoa(i)
		},
	}).Parse(bodyTemplateMap)
	if err != nil {
		return err
	}

	title := "Information about the " + s.cfg.DB.Database + " database, create: " + time.Now().Format("02.01.2006")
	p := &ViewData{Title: title, CreatedOn: time.Now().Format("02.01.2006"), Body: data}
	var body bytes.Buffer
	if err := t.Execute(&body, p); err != nil {
		return err
	}
	bodyMessage := body.String()

	if s.cfg.Mail.IsSendEmail {
		err = email.New(s.cfg).SendText(bodyMessage, title, s.cfg.Mail.ToSendEmail)
		if err != nil {
			return err
		}
	} else {
		fmt.Println(bodyMessage)
	}

	return nil
}
