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
	Id          int64         `json:"id"`
	UserId      sql.NullInt64 `json:"user_id"`
	CandidateId sql.NullInt64 `json:"candidate_id"`
	CompanyId   sql.NullInt64 `json:"company_id"`
	VacancyId   sql.NullInt64 `json:"vacancy_id"`
	Text        string        `json:"text"`
	ActionCode  int64         `json:"action_code"`
	DateCreate  time.Time     `json:"date_create"`
	create_at   time.Time     `json:"create_at"`
}

func (r *LogsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.UserId

	case "candidate_id":
		return &r.CandidateId

	case "company_id":
		return &r.CompanyId

	case "vacancy_id":
		return &r.VacancyId

	case "text":
		return &r.Text

	case "action_code":
		return &r.ActionCode

	case "date_create":
		return &r.DateCreate

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
		return r.UserId

	case "candidate_id":
		return r.CandidateId

	case "company_id":
		return r.CompanyId

	case "vacancy_id":
		return r.VacancyId

	case "text":
		return r.Text

	case "action_code":
		return r.ActionCode

	case "date_create":
		return r.DateCreate

	case "create_at":
		return r.create_at

	default:
		return nil
	}
}

func NewLogs(db *dbEngine.DB) (*Logs, error) {
	table, ok := db.Tables[TABLE_LOGS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_LOGS}
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
