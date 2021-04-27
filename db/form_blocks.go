// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Form_blocks struct {
	dbEngine.Table
	Record *Form_blocksFields
	rows   sql.Rows
}

type Form_blocksFields struct {
	Id          int32       `json:"id"`
	Id_forms    int32       `json:"id_forms"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Tablename   string      `json:"tablename"`
	Columns     []string    `json:"columns"`
	Buttons     interface{} `json:"buttons"`
}

func (r *Form_blocksFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "id_forms":
		return &r.Id_forms

	case "title":
		return &r.Title

	case "description":
		return &r.Description

	case "tablename":
		return &r.Tablename

	case "columns":
		return &r.Columns

	case "buttons":
		return &r.Buttons

	default:
		return nil
	}
}

func (r *Form_blocksFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "id_forms":
		return r.Id_forms

	case "title":
		return r.Title

	case "description":
		return r.Description

	case "tablename":
		return r.Tablename

	case "columns":
		return r.Columns

	case "buttons":
		return r.Buttons

	default:
		return nil
	}
}

func NewForm_blocks(db *dbEngine.DB) (*Form_blocks, error) {
	table, ok := db.Tables[TableFormBlocks]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableFormBlocks}
	}

	return &Form_blocks{
		Table: table,
	}, nil
}

func (t *Form_blocks) NewRecord() *Form_blocksFields {
	t.Record = &Form_blocksFields{}
	return t.Record
}

func (t *Form_blocks) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Form_blocks) SelectSelfScanEach(ctx context.Context, each func(record *Form_blocksFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Form_blocks) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Form_blocks) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
