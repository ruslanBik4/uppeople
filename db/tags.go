// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Tags struct {
	dbEngine.Table
	Record *TagsFields
	rows   sql.Rows
}

type TagsFields struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	Parent_id int64  `json:"parent_id"`
	Order_num int64  `json:"order_num"`
}

func (r *TagsFields) GetFields(columns []dbEngine.Column) []interface{} {
	if len(columns) == 0 {
		return []interface{}{
			&r.Id,
			&r.Name,
			&r.Color,
			&r.Parent_id,
			&r.Order_num,
		}
	}

	v := make([]interface{}, len(columns))
	for i, col := range columns {
		v[i] = r.RefColValue(col.Name())
	}

	return v
}

func (r *TagsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "color":
		return &r.Color

	case "parent_id":
		return &r.Parent_id

	case "order_num":
		return &r.Order_num

	default:
		return nil
	}
}

func (r *TagsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "color":
		return r.Color

	case "parent_id":
		return r.Parent_id

	case "order_num":
		return r.Order_num

	default:
		return nil
	}
}

func NewTags(db *dbEngine.DB) (*Tags, error) {
	table, ok := db.Tables["tags"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "tags"}
	}

	return &Tags{
		Table: table,
	}, nil
}

func (t *Tags) NewRecord() *TagsFields {
	t.Record = &TagsFields{}
	return t.Record
}

func (t *Tags) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Tags) SelectSelfScanEach(ctx context.Context, each func(record *TagsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Tags) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Tags) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
