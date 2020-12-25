// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Status_for_vacs struct {
	dbEngine.Table
	Record *Status_for_vacsFields
	rows   sql.Rows
}

type Status_for_vacsFields struct {
	Id        int64          `json:"id"`
	Status    sql.NullString `json:"status"`
	Color     string         `json:"color"`
	Order_num int64          `json:"order_num"`
}

func (r *Status_for_vacsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "status":
		return &r.Status

	case "color":
		return &r.Color

	case "order_num":
		return &r.Order_num

	default:
		return nil
	}
}

func (r *Status_for_vacsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "status":
		return r.Status

	case "color":
		return r.Color

	case "order_num":
		return r.Order_num

	default:
		return nil
	}
}

func NewStatus_for_vacs(db *dbEngine.DB) (*Status_for_vacs, error) {
	table, ok := db.Tables["status_for_vacs"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "status_for_vacs"}
	}

	return &Status_for_vacs{
		Table: table,
	}, nil
}

func (t *Status_for_vacs) NewRecord() *Status_for_vacsFields {
	t.Record = &Status_for_vacsFields{}
	return t.Record
}

func (t *Status_for_vacs) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Status_for_vacs) SelectSelfScanEach(ctx context.Context, each func(record *Status_for_vacsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Status_for_vacs) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Status_for_vacs) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
