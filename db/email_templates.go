// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Email_templates struct {
	dbEngine.Table
	Record *Email_templatesFields
	rows   sql.Rows
}

type Email_templatesFields struct {
	Id    int64          `json:"id"`
	Theme sql.NullString `json:"theme"`
	Text  sql.NullString `json:"text"`
}

func (r *Email_templatesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "theme":
		return &r.Theme

	case "text":
		return &r.Text

	default:
		return nil
	}
}

func (r *Email_templatesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "theme":
		return r.Theme

	case "text":
		return r.Text

	default:
		return nil
	}
}

func NewEmail_templates(db *dbEngine.DB) (*Email_templates, error) {
	table, ok := db.Tables["email_templates"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "email_templates"}
	}

	return &Email_templates{
		Table: table,
	}, nil
}

func (t *Email_templates) NewRecord() *Email_templatesFields {
	t.Record = &Email_templatesFields{}
	return t.Record
}

func (t *Email_templates) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Email_templates) SelectSelfScanEach(ctx context.Context, each func(record *Email_templatesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Email_templates) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Email_templates) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
