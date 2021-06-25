package api

import (
	"database/sql/driver"
	"strings"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

func getVacToCand(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res db.SelectedUnits) {
	err := DB.Conn.SelectAndScanEach(ctx,
		nil,
		&res,
		`select v.status as id, s.status  as label, lower(s.status) as value
        from vacancies_to_candidates v join status_for_vacs s on (s.id = v.status)
`,
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getLocations(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res db.SelectedUnits) {
	statUses, _ := db.NewLocation_for_vacancies(DB)

	err := statUses.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
		dbEngine.OrderBy("name"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getLanguages(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res db.SelectedUnits) {

	err := DB.Conn.SelectAndScanEach(ctx,
		nil,
		&res,
		"select id, name as label, LOWER(name) as value from languages order by name",
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getRecruiters(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res db.SelectedUnits) {
	users, _ := db.NewUsers(DB)

	err := users.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
		dbEngine.OrderBy("name"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getCompanies(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, opt ...dbEngine.BuildSqlOptions) (res db.SelectedUnits) {
	company, _ := db.NewCompanies(DB)

	err := company.SelectAndScanEach(ctx,
		nil,
		&res,
		append(opt,
			dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
			dbEngine.OrderBy("name"))...,
	)
	if err != nil {
		logs.ErrorLog(err, "	SelectSelfScanEach")
	}

	return
}

func EmptyValue(value interface{}) bool {
	if value == nil {
		return true
	}

	switch val := value.(type) {
	case []int32:
		return len(val) == 0
	case []int64:
		return len(val) == 0
	case []float32:
		return len(val) == 0
	case []float64:
		return len(val) == 0
	case []string:
		return len(val) == 0
	case int32:
		return val == 0
	case int64:
		return val == 0
	case float32:
		return val == 0
	case float64:
		return val == 0
	case time.Time:
		return val.IsZero()
	case *time.Time:
		if val == nil {
			return true
		} else {
			return val.IsZero()
		}
	case driver.Valuer:
		v1, _ := val.Value()
		return v1 == nil
	case string:
		return strings.TrimSpace(val) == ""
	default:
		return false
	}
}
