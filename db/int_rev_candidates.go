// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type IntRevCandidates struct {
	dbEngine.Table
	Record *IntRevCandidatesFields
	rows   sql.Rows
}

type IntRevCandidatesFields struct {
	CandidateId sql.NullInt64 `json:"candidate_id"`
	CompanyId   sql.NullInt64 `json:"company_id"`
	VacancyId   sql.NullInt64 `json:"vacancy_id"`
	Status      sql.NullInt64 `json:"status"`
	UserId      sql.NullInt64 `json:"user_id"`
	Date        time.Time     `json:"date"`
}

func (r *IntRevCandidatesFields) RefColValue(name string) interface{} {
	switch name {
	case "candidate_id":
		return &r.CandidateId

	case "company_id":
		return &r.CompanyId

	case "vacancy_id":
		return &r.VacancyId

	case "status":
		return &r.Status

	case "user_id":
		return &r.UserId

	case "date":
		return &r.Date

	default:
		return nil
	}
}

func (r *IntRevCandidatesFields) ColValue(name string) interface{} {
	switch name {
	case "candidate_id":
		return r.CandidateId

	case "company_id":
		return r.CompanyId

	case "vacancy_id":
		return r.VacancyId

	case "status":
		return r.Status

	case "user_id":
		return r.UserId

	case "date":
		return r.Date

	default:
		return nil
	}
}

func NewInt_rev_candidates(db *dbEngine.DB) (*IntRevCandidates, error) {
	table, ok := db.Tables[TABLE_INT_REV_CANDIDATES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_INT_REV_CANDIDATES}
	}

	return &IntRevCandidates{
		Table: table,
	}, nil
}

func (t *IntRevCandidates) NewRecord() *IntRevCandidatesFields {
	t.Record = &IntRevCandidatesFields{}
	return t.Record
}

func (t *IntRevCandidates) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *IntRevCandidates) SelectSelfScanEach(ctx context.Context, each func(record *IntRevCandidatesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *IntRevCandidates) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *IntRevCandidates) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
