package controller

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"kvart-info/internal/model"
	"kvart-info/pkg/utils"
	"strconv"
	"text/tabwriter"
	"time"
)

// OutputInfo вывод информации или на почту или на экран
func (c *Controller) OutputInfo() (string, string, error) {

	data, err := c.uc.GetSummaryInfo(context.Background())
	if err != nil {
		return "", "", err
	}

	bodyMessage, title, err := c.CreateBodyText(data)
	if err != nil {
		return "", "", err
	}

	return bodyMessage, title, nil
}

// CreateBodyText формируем письмо
func (c *Controller) CreateBodyText(data []model.SummaryInfo) (string, string, error) {

	type ViewData struct {
		Title     string
		CreatedOn string
		Body      []model.SummaryInfo
	}
	title := fmt.Sprintf("Information about the %s database, create: %s",
		c.cfg.DB.Database,
		time.Now().Format("02.01.2006"),
	)
	var body bytes.Buffer

	if c.cfg.IsSendEmail {
		t, err := template.New("").Funcs(template.FuncMap{
			"IncInt": func(i int) string {
				i++
				return strconv.Itoa(i)
			},
			"getValidName": func(src string) string {
				return utils.GetValidFileName(src)
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
	t := template.Must(template.New("").Funcs(template.FuncMap{
		"IncInt": func(i int) string {
			i++
			return strconv.Itoa(i)
		},
		"getValidName": func(src string) string {
			return utils.GetValidFileName(src)
		},
	}).Parse(tmplText))
	w := tabwriter.NewWriter(&body, 5, 0, 3, ' ', 0)
	if err := t.Execute(w, data); err != nil {
		return "", "", err
	}
	_ = w.Flush()
	bodyMessage := body.String()

	return bodyMessage, title, nil
}

var bodyTemplateMap = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>

    <style>
        table, th, td {
            border: 1px solid black;
            border-collapse: collapse;
        }
        table {
            margin-bottom: 50px;
        }
		td.val {
			text-align: right;
		}
    </style>

</head>
<body>
<table class="table">
	<th>№</th>
    <th>Фин.период</th>
    <th>Тип фонда</th>
    <th>Домов</th>
    <th>Лицевых</th>
    <th>Разница с пред.</th>
    <th>Помещений</th>
    <th>Разница с пред.</th>
    <th>Площадь</th>
    <th>Разница с пред.</th>
    {{ $CreatedOn := .CreatedOn}}
    {{range $i, $el := .Body}}
    <tr>
        <td>
		{{if ne .TipName "Итого"}}
		{{(IncInt $i)}}
		{{end}}
        </td>
        <td>
        {{if ne .TipName "Итого"}}
        {{ .FinPeriod }} {{else}} {{ $CreatedOn }} 
        {{end}}
        </td>
        <td>{{(getValidName .TipName)}}</td>	
		<td class="val">{{.CountBuild}}</td>
        <td class="val">{{.CountLic}}</td>
        <td class="val">{{.KolOccDif}}</td>
        <td class="val">{{.Flat}}</td>
        <td class="val">{{.KolFlatDif}}</td>
        <td class="val">{{printf "%.2f" .TotalSq}}</td>
        <td class="val">{{.TotalSqDif}}</td>
    </tr>
    {{end}}
</table>
<hr>
Данное электронное письмо было направлено вам автоматизированной системой.<br>
Не отвечайте по данному адресу электронной почты.
</body>
</html>
`

const tmplText = "\n" +
	"Фин.период\tТип фонда\tДомов\tЛицевых\tПомещений\tПлощадь\tРазница лицевых\n" +
	"{{range .}}{{.FinPeriod}}\t{{(getValidName .TipName)}}\t{{.CountBuild}}\t{{.CountLic}}\t{{.Flat}}\t{{.TotalSq}}\t{{.KolOccDif}}\n{{end}}"
