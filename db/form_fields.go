// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Form_fields struct {
	dbEngine.Table
	Record *Form_fieldsFields
	rows   sql.Rows
}

type Form_fieldsFields struct {
	Id             int32  `json:"id"`
	Name           string `json:"name"`
	Input_type     string `json:"input_type"`
	Title          string `json:"title"`
	Pattern        string `json:"pattern"`
	Placeholder    string `json:"placeholder"`
	Autofocus      bool   `json:"autofocus"`
	Disabled       bool   `json:"disabled"`
	Required       bool   `json:"required"`
	Readonly       bool   `json:"readonly"`
	Id_form_blocks int32  `json:"id_form_blocks"`
}

func (r *Form_fieldsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "input_type":
		return &r.Input_type

	case "title":
		return &r.Title

	case "pattern":
		return &r.Pattern

	case "placeholder":
		return &r.Placeholder

	case "autofocus":
		return &r.Autofocus

	case "disabled":
		return &r.Disabled

	case "required":
		return &r.Required

	case "readonly":
		return &r.Readonly

	case "id_form_blocks":
		return &r.Id_form_blocks

	default:
		return nil
	}
}

func (r *Form_fieldsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "input_type":
		return r.Input_type

	case "title":
		return r.Title

	case "pattern":
		return r.Pattern

	case "placeholder":
		return r.Placeholder

	case "autofocus":
		return r.Autofocus

	case "disabled":
		return r.Disabled

	case "required":
		return r.Required

	case "readonly":
		return r.Readonly

	case "id_form_blocks":
		return r.Id_form_blocks

	default:
		return nil
	}
}

func NewForm_fields(db *dbEngine.DB) (*Form_fields, error) {
	table, ok := db.Tables["form_fields"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "form_fields"}
	}

	return &Form_fields{
		Table: table,
	}, nil
}

func (t *Form_fields) NewRecord() *Form_fieldsFields {
	t.Record = &Form_fieldsFields{}
	return t.Record
}

func (t *Form_fields) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Form_fields) SelectSelfScanEach(ctx context.Context, each func(record *Form_fieldsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Form_fields) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Form_fields) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
