// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Vacancies_to_candidates struct {
	dbEngine.Table
	Record *Vacancies_to_candidatesFields
	rows   sql.Rows
}

type Vacancies_to_candidatesFields struct {
	Id               int64          `json:"id"`
	Candidate_id     sql.NullInt64  `json:"candidate_id"`
	Company_id       sql.NullInt64  `json:"company_id"`
	Vacancy_id       sql.NullInt64  `json:"vacancy_id"`
	Status           sql.NullInt64  `json:"status"`
	User_id          sql.NullInt64  `json:"user_id"`
	Date_create      time.Time      `json:"date_create"`
	Date_last_change time.Time      `json:"date_last_change"`
	Rej_text         sql.NullString `json:"rej_text"`
	Rating           sql.NullInt64  `json:"rating"`
	Notice           sql.NullString `json:"notice"`
}

func (r *Vacancies_to_candidatesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "candidate_id":
		return &r.Candidate_id

	case "company_id":
		return &r.Company_id

	case "vacancy_id":
		return &r.Vacancy_id

	case "status":
		return &r.Status

	case "user_id":
		return &r.User_id

	case "date_create":
		return &r.Date_create

	case "date_last_change":
		return &r.Date_last_change

	case "rej_text":
		return &r.Rej_text

	case "rating":
		return &r.Rating

	case "notice":
		return &r.Notice

	default:
		return nil
	}
}

func (r *Vacancies_to_candidatesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "candidate_id":
		return r.Candidate_id

	case "company_id":
		return r.Company_id

	case "vacancy_id":
		return r.Vacancy_id

	case "status":
		return r.Status

	case "user_id":
		return r.User_id

	case "date_create":
		return r.Date_create

	case "date_last_change":
		return r.Date_last_change

	case "rej_text":
		return r.Rej_text

	case "rating":
		return r.Rating

	case "notice":
		return r.Notice

	default:
		return nil
	}
}

func NewVacancies_to_candidates(db *dbEngine.DB) (*Vacancies_to_candidates, error) {
	table, ok := db.Tables["vacancies_to_candidates"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "vacancies_to_candidates"}
	}

	return &Vacancies_to_candidates{
		Table: table,
	}, nil
}

func (t *Vacancies_to_candidates) NewRecord() *Vacancies_to_candidatesFields {
	t.Record = &Vacancies_to_candidatesFields{}
	return t.Record
}

func (t *Vacancies_to_candidates) GetFields(columns []dbEngine.Column) []interface{} {
	if len(columns) == 0 {
		columns = t.Columns()
	}

	t.NewRecord()
	v := make([]interface{}, len(columns))
	for i, col := range columns {
		v[i] = t.Record.RefColValue(col.Name())
	}

	return v
}

func (t *Vacancies_to_candidates) SelectSelfScanEach(ctx context.Context, each func(record *Vacancies_to_candidatesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Vacancies_to_candidates) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
	if len(Options) == 0 {
		v := make([]interface{}, len(t.Columns()))
		columns := make([]string, len(t.Columns()))
		for i, col := range t.Columns() {
			columns[i] = col.Name()
			v[i] = t.Record.ColValue(col.Name())
		}
		Options = append(Options,
			dbEngine.ColumnsForSelect(columns...),
			dbEngine.ArgsForSelect(v...))
	}

	return t.Table.Insert(ctx, Options...)
}

func (t *Vacancies_to_candidates) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
	if len(Options) == 0 {
		v := make([]interface{}, len(t.Columns()))
		priV := make([]interface{}, 0)
		columns := make([]string, 0, len(t.Columns()))
		priColumns := make([]string, 0, len(t.Columns()))
		for _, col := range t.Columns() {
			if col.Primary() {
				priColumns = append(priColumns, col.Name())
				priV[len(priColumns)-1] = t.Record.ColValue(col.Name())
				continue
			}

			columns = append(columns, col.Name())
			v[len(columns)-1] = t.Record.ColValue(col.Name())
		}

		Options = append(
			Options,
			dbEngine.ColumnsForSelect(columns...),
			dbEngine.WhereForSelect(priColumns...),
			dbEngine.ArgsForSelect(append(v, priV...)...),
		)
	}

	return t.Table.Update(ctx, Options...)
}
