// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Int_rev_candidates struct {
	dbEngine.Table
	Record *Int_rev_candidatesFields
	rows   sql.Rows
}

type Int_rev_candidatesFields struct {
	Id           int64         `json:"id"`
	Candidate_id sql.NullInt64 `json:"candidate_id"`
	Company_id   sql.NullInt64 `json:"company_id"`
	Vacancy_id   sql.NullInt64 `json:"vacancy_id"`
	Status       sql.NullInt64 `json:"status"`
	User_id      sql.NullInt64 `json:"user_id"`
	Date         time.Time     `json:"date"`
}

func (r *Int_rev_candidatesFields) RefColValue(name string) interface{} {
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

	case "date":
		return &r.Date

	default:
		return nil
	}
}

func (r *Int_rev_candidatesFields) ColValue(name string) interface{} {
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

	case "date":
		return r.Date

	default:
		return nil
	}
}

func NewInt_rev_candidates(db *dbEngine.DB) (*Int_rev_candidates, error) {
	table, ok := db.Tables[TABLE_IntRevCandidates]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_IntRevCandidates}
	}

	return &Int_rev_candidates{
		Table: table,
	}, nil
}

func (t *Int_rev_candidates) NewRecord() *Int_rev_candidatesFields {
	t.Record = &Int_rev_candidatesFields{}
	return t.Record
}

func (t *Int_rev_candidates) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Int_rev_candidates) SelectSelfScanEach(ctx context.Context, each func(record *Int_rev_candidatesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Int_rev_candidates) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Int_rev_candidates) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
