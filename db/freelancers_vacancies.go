// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Freelancers_vacancies struct {
	dbEngine.Table
	Record *Freelancers_vacanciesFields
	rows   sql.Rows
}

type Freelancers_vacanciesFields struct {
	Id            int64         `json:"id"`
	Vacancy_id    sql.NullInt64 `json:"vacancy_id"`
	User_id       sql.NullInt64 `json:"user_id"`
	Freelancer_id sql.NullInt64 `json:"freelancer_id"`
	Date          time.Time     `json:"date"`
}

func (r *Freelancers_vacanciesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "vacancy_id":
		return &r.Vacancy_id

	case "user_id":
		return &r.User_id

	case "freelancer_id":
		return &r.Freelancer_id

	case "date":
		return &r.Date

	default:
		return nil
	}
}

func (r *Freelancers_vacanciesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "vacancy_id":
		return r.Vacancy_id

	case "user_id":
		return r.User_id

	case "freelancer_id":
		return r.Freelancer_id

	case "date":
		return r.Date

	default:
		return nil
	}
}

func NewFreelancers_vacancies(db *dbEngine.DB) (*Freelancers_vacancies, error) {
	table, ok := db.Tables["freelancers_vacancies"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "freelancers_vacancies"}
	}

	return &Freelancers_vacancies{
		Table: table,
	}, nil
}

func (t *Freelancers_vacancies) NewRecord() *Freelancers_vacanciesFields {
	t.Record = &Freelancers_vacanciesFields{}
	return t.Record
}

func (t *Freelancers_vacancies) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Freelancers_vacancies) SelectSelfScanEach(ctx context.Context, each func(record *Freelancers_vacanciesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Freelancers_vacancies) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Freelancers_vacancies) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
