// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Management struct {
	dbEngine.Table
	Record *ManagementFields
	rows   sql.Rows
}

type ManagementFields struct {
	Id_users int32       `json:"id_users"`
	Id_firma int32       `json:"id_firma"`
	Grants   interface{} `json:"grants"`
}

func (r *ManagementFields) RefColValue(name string) interface{} {
	switch name {
	case "id_users":
		return &r.Id_users

	case "id_firma":
		return &r.Id_firma

	case "grants":
		return &r.Grants

	default:
		return nil
	}
}

func (r *ManagementFields) ColValue(name string) interface{} {
	switch name {
	case "id_users":
		return r.Id_users

	case "id_firma":
		return r.Id_firma

	case "grants":
		return r.Grants

	default:
		return nil
	}
}

func NewManagement(db *dbEngine.DB) (*Management, error) {
	table, ok := db.Tables[TABLE_MANAGEMENT]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_MANAGEMENT}
	}

	return &Management{
		Table: table,
	}, nil
}

func (t *Management) NewRecord() *ManagementFields {
	t.Record = &ManagementFields{}
	return t.Record
}

func (t *Management) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Management) SelectSelfScanEach(ctx context.Context, each func(record *ManagementFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Management) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Management) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
