package services

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
	"Фин.период\tТип фонда\tДомов\tЛицевых\tПомещений\tПлощадь\n" +
	"{{range .}}{{.FinPeriod}}\t{{(getValidName .TipName)}}\t{{.CountBuild}}\t{{.CountLic}}\t{{.Flat}}\t{{.TotalSq}}\n{{end}}"
