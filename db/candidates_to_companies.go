// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Candidates_to_companies struct {
	dbEngine.Table
	Record *Candidates_to_companiesFields
	rows   sql.Rows
}

type Candidates_to_companiesFields struct {
	Id           int64         `json:"id"`
	Candidate_id sql.NullInt64 `json:"candidate_id"`
	Company_id   sql.NullInt64 `json:"company_id"`
	Visible      sql.NullInt64 `json:"visible"`
}

func (r *Candidates_to_companiesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "candidate_id":
		return &r.Candidate_id

	case "company_id":
		return &r.Company_id

	case "visible":
		return &r.Visible

	default:
		return nil
	}
}

func (r *Candidates_to_companiesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "candidate_id":
		return r.Candidate_id

	case "company_id":
		return r.Company_id

	case "visible":
		return r.Visible

	default:
		return nil
	}
}

func NewCandidates_to_companies(db *dbEngine.DB) (*Candidates_to_companies, error) {
	table, ok := db.Tables[TABLE_CandidatesToCompanies]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_CandidatesToCompanies}
	}

	return &Candidates_to_companies{
		Table: table,
	}, nil
}

func (t *Candidates_to_companies) NewRecord() *Candidates_to_companiesFields {
	t.Record = &Candidates_to_companiesFields{}
	return t.Record
}

func (t *Candidates_to_companies) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Candidates_to_companies) SelectSelfScanEach(ctx context.Context, each func(record *Candidates_to_companiesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Candidates_to_companies) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Candidates_to_companies) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
