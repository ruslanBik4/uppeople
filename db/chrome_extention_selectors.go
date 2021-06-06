// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Chrome_extention_selectors struct {
	dbEngine.Table
	Record *Chrome_extention_selectorsFields
	rows   sql.Rows
}

type Chrome_extention_selectorsFields struct {
	Id       int64  `json:"id"`
	Param    string `json:"param"`
	Selector string `json:"selector"`
	Type     string `json:"type"`
}

func (r *Chrome_extention_selectorsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "param":
		return &r.Param

	case "selector":
		return &r.Selector

	case "type":
		return &r.Type

	default:
		return nil
	}
}

func (r *Chrome_extention_selectorsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "param":
		return r.Param

	case "selector":
		return r.Selector

	case "type":
		return r.Type

	default:
		return nil
	}
}

func NewChrome_extention_selectors(db *dbEngine.DB) (*Chrome_extention_selectors, error) {
	table, ok := db.Tables[TABLE_CHROME_EXTENTION_SELECTORS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_CHROME_EXTENTION_SELECTORS}
	}

	return &Chrome_extention_selectors{
		Table: table,
	}, nil
}

func (t *Chrome_extention_selectors) NewRecord() *Chrome_extention_selectorsFields {
	t.Record = &Chrome_extention_selectorsFields{}
	return t.Record
}

func (t *Chrome_extention_selectors) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Chrome_extention_selectors) SelectSelfScanEach(ctx context.Context, each func(record *Chrome_extention_selectorsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Chrome_extention_selectors) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Chrome_extention_selectors) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
