package model

// TotalData ...
type TotalData struct {
	FinID      int64   `db:"fin_id"`
	FinPeriod  string  `db:"fin_period"`
	TipName    string  `db:"tip_name"`
	CountBuild int64   `db:"count_build"`
	CountLic   int64   `db:"count_occ"`
	Flat       int64   `db:"flat"`
	TotalSq    float64 `db:"total_sq"`
	KolOccDif  int64   `db:"kol_occ_dif"`
	KolFlatDif int64   `db:"kol_flat_dif"`
	TotalSqDif float64 `db:"total_sq_dif"`
}
