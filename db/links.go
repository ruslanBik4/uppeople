// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Links struct {
	dbEngine.Table
	Record *LinksFields
	rows   sql.Rows
}

type LinksFields struct {
	Id    int64          `json:"id"`
	Title sql.NullString `json:"title"`
	Link  sql.NullString `json:"link"`
}

func (r *LinksFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "title":
		return &r.Title

	case "link":
		return &r.Link

	default:
		return nil
	}
}

func (r *LinksFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "title":
		return r.Title

	case "link":
		return r.Link

	default:
		return nil
	}
}

func NewLinks(db *dbEngine.DB) (*Links, error) {
	table, ok := db.Tables["links"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "links"}
	}

	return &Links{
		Table: table,
	}, nil
}

func (t *Links) NewRecord() *LinksFields {
	t.Record = &LinksFields{}
	return t.Record
}

func (t *Links) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Links) SelectSelfScanEach(ctx context.Context, each func(record *LinksFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Links) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Links) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
