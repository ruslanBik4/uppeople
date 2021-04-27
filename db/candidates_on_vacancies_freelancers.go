// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Candidates_on_vacancies_freelancers struct {
	dbEngine.Table
	Record *Candidates_on_vacancies_freelancersFields
	rows   sql.Rows
}

type Candidates_on_vacancies_freelancersFields struct {
	Id            int64         `json:"id"`
	Candidate_id  sql.NullInt64 `json:"candidate_id"`
	Company_id    sql.NullInt64 `json:"company_id"`
	Vacancy_id    sql.NullInt64 `json:"vacancy_id"`
	Freelancer_id sql.NullInt64 `json:"freelancer_id"`
	Date          time.Time     `json:"date"`
}

func (r *Candidates_on_vacancies_freelancersFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "candidate_id":
		return &r.Candidate_id

	case "company_id":
		return &r.Company_id

	case "vacancy_id":
		return &r.Vacancy_id

	case "freelancer_id":
		return &r.Freelancer_id

	case "date":
		return &r.Date

	default:
		return nil
	}
}

func (r *Candidates_on_vacancies_freelancersFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "candidate_id":
		return r.Candidate_id

	case "company_id":
		return r.Company_id

	case "vacancy_id":
		return r.Vacancy_id

	case "freelancer_id":
		return r.Freelancer_id

	case "date":
		return r.Date

	default:
		return nil
	}
}

func NewCandidates_on_vacancies_freelancers(db *dbEngine.DB) (*Candidates_on_vacancies_freelancers, error) {
	table, ok := db.Tables[TableCandidatesOnVacanciesFreelancers]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableCandidatesOnVacanciesFreelancers}
	}

	return &Candidates_on_vacancies_freelancers{
		Table: table,
	}, nil
}

func (t *Candidates_on_vacancies_freelancers) NewRecord() *Candidates_on_vacancies_freelancersFields {
	t.Record = &Candidates_on_vacancies_freelancersFields{}
	return t.Record
}

func (t *Candidates_on_vacancies_freelancers) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Candidates_on_vacancies_freelancers) SelectSelfScanEach(ctx context.Context, each func(record *Candidates_on_vacancies_freelancersFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Candidates_on_vacancies_freelancers) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Candidates_on_vacancies_freelancers) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
