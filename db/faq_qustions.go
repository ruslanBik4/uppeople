// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Faq_qustions struct {
	dbEngine.Table
	Record *Faq_qustionsFields
	rows   sql.Rows
}

type Faq_qustionsFields struct {
	Id          int64          `json:"id"`
	Cat_id      sql.NullInt64  `json:"cat_id"`
	Question    sql.NullString `json:"question"`
	Answer      sql.NullString `json:"answer"`
	Date_create time.Time      `json:"date_create"`
}

func (r *Faq_qustionsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "cat_id":
		return &r.Cat_id

	case "question":
		return &r.Question

	case "answer":
		return &r.Answer

	case "date_create":
		return &r.Date_create

	default:
		return nil
	}
}

func (r *Faq_qustionsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "cat_id":
		return r.Cat_id

	case "question":
		return r.Question

	case "answer":
		return r.Answer

	case "date_create":
		return r.Date_create

	default:
		return nil
	}
}

func NewFaq_qustions(db *dbEngine.DB) (*Faq_qustions, error) {
	table, ok := db.Tables[TableFAQQuestions]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableFAQQuestions}
	}

	return &Faq_qustions{
		Table: table,
	}, nil
}

func (t *Faq_qustions) NewRecord() *Faq_qustionsFields {
	t.Record = &Faq_qustionsFields{}
	return t.Record
}

func (t *Faq_qustions) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Faq_qustions) SelectSelfScanEach(ctx context.Context, each func(record *Faq_qustionsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Faq_qustions) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Faq_qustions) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
