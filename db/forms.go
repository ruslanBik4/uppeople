// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Forms struct {
	dbEngine.Table
	Record *FormsFields
	rows   sql.Rows
}

type FormsFields struct {
	Id          int32       `json:"id"`
	Title       string      `json:"title"`
	Action      string      `json:"action"`
	Post        bool        `json:"post"`
	Description string      `json:"description"`
	HideBlock   interface{} `json:"hideblock"`
}

func (r *FormsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "title":
		return &r.Title

	case "action":
		return &r.Action

	case "post":
		return &r.Post

	case "description":
		return &r.Description

	case "hideblock":
		return &r.HideBlock

	default:
		return nil
	}
}

func (r *FormsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "title":
		return r.Title

	case "action":
		return r.Action

	case "post":
		return r.Post

	case "description":
		return r.Description

	case "hideblock":
		return r.HideBlock

	default:
		return nil
	}
}

func NewForms(db *dbEngine.DB) (*Forms, error) {
	table, ok := db.Tables[TableForms]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableForms}
	}

	return &Forms{
		Table: table,
	}, nil
}

func (t *Forms) NewRecord() *FormsFields {
	t.Record = &FormsFields{}
	return t.Record
}

func (t *Forms) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Forms) SelectSelfScanEach(ctx context.Context, each func(record *FormsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Forms) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Forms) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
