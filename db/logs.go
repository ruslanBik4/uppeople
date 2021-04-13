// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Logs struct {
	dbEngine.Table
	Record *LogsFields
	rows   sql.Rows
}

type LogsFields struct {
	Id            int64         `json:"id"`
	User_id       sql.NullInt64 `json:"user_id"`
	Candidate_id  sql.NullInt64 `json:"candidate_id"`
	Company_id    sql.NullInt64 `json:"company_id"`
	Vacancy_id    sql.NullInt64 `json:"vacancy_id"`
	Text          string        `json:"text"`
	Kod_deystviya int64         `json:"kod_deystviya"`
	Date_create   time.Time     `json:"date_create"`
	create_at     time.Time     `json:"create_at"`
}

func (r *LogsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "candidate_id":
		return &r.Candidate_id

	case "company_id":
		return &r.Company_id

	case "vacancy_id":
		return &r.Vacancy_id

	case "text":
		return &r.Text

	case "kod_deystviya":
		return &r.Kod_deystviya

	case "date_create":
		return &r.Date_create

	case "create_at":
		return &r.create_at

	default:
		return nil
	}
}

func (r *LogsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "candidate_id":
		return r.Candidate_id

	case "company_id":
		return r.Company_id

	case "vacancy_id":
		return r.Vacancy_id

	case "text":
		return r.Text

	case "kod_deystviya":
		return r.Kod_deystviya

	case "date_create":
		return r.Date_create

	case "create_at":
		return r.create_at

	default:
		return nil
	}
}

func NewLogs(db *dbEngine.DB) (*Logs, error) {
	table, ok := db.Tables["logs"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "logs"}
	}

	return &Logs{
		Table: table,
	}, nil
}

func (t *Logs) NewRecord() *LogsFields {
	t.Record = &LogsFields{}
	return t.Record
}

func (t *Logs) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Logs) SelectSelfScanEach(ctx context.Context, each func(record *LogsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Logs) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Logs) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
