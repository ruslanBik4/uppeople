// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Meetings struct {
	dbEngine.Table
	Record *MeetingsFields
	rows   sql.Rows
}

type MeetingsFields struct {
	Id           int64          `json:"id"`
	Candidate_id sql.NullInt64  `json:"candidate_id"`
	Company_id   sql.NullInt64  `json:"company_id"`
	Vacancy_id   sql.NullInt64  `json:"vacancy_id"`
	Title        sql.NullString `json:"title"`
	Date         time.Time      `json:"date"`
	D_t          time.Time      `json:"d_t"`
	User_id      sql.NullInt64  `json:"user_id"`
	Color        sql.NullString `json:"color"`
	Type         sql.NullString `json:"type"`
}

func (r *MeetingsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "candidate_id":
		return &r.Candidate_id

	case "company_id":
		return &r.Company_id

	case "vacancy_id":
		return &r.Vacancy_id

	case "title":
		return &r.Title

	case "date":
		return &r.Date

	case "d_t":
		return &r.D_t

	case "user_id":
		return &r.User_id

	case "color":
		return &r.Color

	case "type":
		return &r.Type

	default:
		return nil
	}
}

func (r *MeetingsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "candidate_id":
		return r.Candidate_id

	case "company_id":
		return r.Company_id

	case "vacancy_id":
		return r.Vacancy_id

	case "title":
		return r.Title

	case "date":
		return r.Date

	case "d_t":
		return r.D_t

	case "user_id":
		return r.User_id

	case "color":
		return r.Color

	case "type":
		return r.Type

	default:
		return nil
	}
}

func NewMeetings(db *dbEngine.DB) (*Meetings, error) {
	table, ok := db.Tables[TableMeetings]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableMeetings}
	}

	return &Meetings{
		Table: table,
	}, nil
}

func (t *Meetings) NewRecord() *MeetingsFields {
	t.Record = &MeetingsFields{}
	return t.Record
}

func (t *Meetings) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Meetings) SelectSelfScanEach(ctx context.Context, each func(record *MeetingsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Meetings) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Meetings) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
