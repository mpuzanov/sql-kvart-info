package repository

import (
	"context"
	"database/sql"
	"kvart-info/internal/model"
	"kvart-info/pkg/mssql"

	"github.com/pkg/errors"
)

// Repository ...
type Repository struct {
	*mssql.MSSQL
}

// New ...
func New(ms *mssql.MSSQL) *Repository {
	return &Repository{ms}
}

// Get получаем сводную информацию из БД
func (r *Repository) Get(ctx context.Context) ([]model.SummaryInfo, error) {

	var data []model.SummaryInfo
	stmt, err := r.DB.PrepareNamedContext(ctx, QuerySummaryInfo)
	if err != nil {
		return nil, errors.Wrap(err, "Get PrepareNamedContext")
	}
	err = stmt.SelectContext(ctx, &data, map[string]interface{}{})
	if err == sql.ErrNoRows {
		return data, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "Get SelectContext")
	}

	return data, nil
}

// QuerySummaryInfo Запрос для получения итогов по БД
var QuerySummaryInfo = `
SELECT max(ot.fin_id) AS fin_id
     , FORMAT(max(t_cur.start_date), 'MMMM yyyy') AS fin_period
     , coalesce(ot.name,N'Итого') AS tip_name
     , sum(t_cur.count_build) AS count_build
     , sum(t_cur.kol_occ) AS count_occ
     , sum(t_cur.kol_flat) AS flat
     , sum(t_cur.total_sq) AS total_sq
     , sum(t_cur.kol_occ - t_prev.kol_occ)  as kol_occ_dif 
     , sum(t_cur.kol_flat - t_prev.kol_flat) as kol_flat_dif
     , sum(t_cur.total_sq - t_prev.total_sq) as total_sq_dif
FROM dbo.Occupation_Types ot
CROSS APPLY  (
         SELECT o.fin_id
              , MAX(o.start_date) AS start_date
              , o.tip_id
              , COUNT(distinct  o.build_id) as count_build
              , COUNT(o.occ) AS kol_occ
              , COUNT(DISTINCT o.flat_id) AS kol_flat
              , SUM(o.total_sq) AS total_sq
         FROM dbo.View_occ_all_lite o
                  JOIN dbo.View_build_all_lite as b ON o.build_id=b.build_id and o.fin_id=b.fin_id
         WHERE o.status_id <> N'закр'
           AND o.tip_id=ot.id
           AND o.fin_id=ot.fin_id
           AND b.is_paym_build=1
           AND o.total_sq > 0
           AND (o.PaidAll<>0 OR o.PaymAccount<>0)
         GROUP BY o.fin_id, o.tip_id
     ) AS t_cur
CROSS APPLY  (
         SELECT o.fin_id
              , MAX(o.start_date) AS start_date
              , o.tip_id
              , COUNT(distinct  o.build_id) as count_build
              , COUNT(o.occ) AS kol_occ
              , COUNT(DISTINCT o.flat_id) AS kol_flat
              , SUM(o.total_sq) AS total_sq
         FROM dbo.View_occ_all_lite o
                  JOIN dbo.View_build_all_lite as b ON o.build_id=b.build_id and o.fin_id=b.fin_id
         WHERE o.status_id <> N'закр'
           AND o.tip_id=ot.id
           AND o.fin_id=ot.fin_id-1
           AND b.is_paym_build=1
           AND o.total_sq > 0
           AND (o.PaidAll<>0 OR o.PaymAccount<>0)
         GROUP BY o.fin_id, o.tip_id
     ) AS t_prev
WHERE ot.payms_value=1 
	and ot.raschet_no=0
GROUP BY ot.name WITH ROLLUP
`
