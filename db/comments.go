// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Comments struct {
	dbEngine.Table
	Record *CommentsFields
	rows   sql.Rows
}

type CommentsFields struct {
	Id           int64          `json:"id"`
	User_id      sql.NullInt64  `json:"user_id"`
	Candidate_id sql.NullInt64  `json:"candidate_id"`
	Company_id   sql.NullInt64  `json:"company_id"`
	Name         sql.NullString `json:"name"`
	Text_comment sql.NullString `json:"text_comment"`
	Date_create  time.Time      `json:"date_create"`
}

func (r *CommentsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "candidate_id":
		return &r.Candidate_id

	case "company_id":
		return &r.Company_id

	case "name":
		return &r.Name

	case "text_comment":
		return &r.Text_comment

	case "date_create":
		return &r.Date_create

	default:
		return nil
	}
}

func (r *CommentsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "candidate_id":
		return r.Candidate_id

	case "company_id":
		return r.Company_id

	case "name":
		return r.Name

	case "text_comment":
		return r.Text_comment

	case "date_create":
		return r.Date_create

	default:
		return nil
	}
}

func NewComments(db *dbEngine.DB) (*Comments, error) {
	table, ok := db.Tables[TableComments]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableComments}
	}

	return &Comments{
		Table: table,
	}, nil
}

func (t *Comments) NewRecord() *CommentsFields {
	t.Record = &CommentsFields{}
	return t.Record
}

func (t *Comments) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Comments) SelectSelfScanEach(ctx context.Context, each func(record *CommentsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Comments) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Comments) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
