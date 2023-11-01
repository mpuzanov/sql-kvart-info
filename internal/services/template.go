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
		{{if ne .tip_name "Итого"}}
		{{(IncInt $i)}}
		{{end}}
        </td>
        <td>
        {{if ne .tip_name "Итого"}}
        {{ .fin_period }} {{else}} {{ $CreatedOn }} 
        {{end}}
        </td>
        <td>{{.tip_name}}</td>	
		<td class="val">{{.count_build}}</td>
        <td class="val">{{.count_occ}}</td>
        <td class="val">{{.kol_occ_dif}}</td>
        <td class="val">{{.flat}}</td>
        <td class="val">{{.kol_flat_dif}}</td>
        <td class="val">{{printf "%.2f" .total_sq}}</td>
        <td class="val">{{.total_sq_dif}}</td>
    </tr>
    {{end}}
</table>
<hr>
Данное электронное письмо было направлено вам автоматизированной системой.<br>
Не отвечайте по данному адресу электронной почты.
</body>
</html>
`
