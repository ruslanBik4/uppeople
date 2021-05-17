// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_options struct {
	dbEngine.Table
	Record *Wp_optionsFields
	rows   sql.Rows
}

type Wp_optionsFields struct {
	Option_id    int64  `json:"option_id"`
	Option_name  string `json:"option_name"`
	Option_value string `json:"option_value"`
	Autoload     string `json:"autoload"`
}

func (r *Wp_optionsFields) RefColValue(name string) interface{} {
	switch name {
	case "option_id":
		return &r.Option_id

	case "option_name":
		return &r.Option_name

	case "option_value":
		return &r.Option_value

	case "autoload":
		return &r.Autoload

	default:
		return nil
	}
}

func (r *Wp_optionsFields) ColValue(name string) interface{} {
	switch name {
	case "option_id":
		return r.Option_id

	case "option_name":
		return r.Option_name

	case "option_value":
		return r.Option_value

	case "autoload":
		return r.Autoload

	default:
		return nil
	}
}

func NewWp_options(db *dbEngine.DB) (*Wp_options, error) {
	table, ok := db.Tables[TABLE_WP_OPTIONS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_WP_OPTIONS}
	}

	return &Wp_options{
		Table: table,
	}, nil
}

func (t *Wp_options) NewRecord() *Wp_optionsFields {
	t.Record = &Wp_optionsFields{}
	return t.Record
}

func (t *Wp_options) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_options) SelectSelfScanEach(ctx context.Context, each func(record *Wp_optionsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_options) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_options) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
