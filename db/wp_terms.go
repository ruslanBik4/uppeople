// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_terms struct {
	dbEngine.Table
	Record *Wp_termsFields
	rows   sql.Rows
}

type Wp_termsFields struct {
	Term_id    int64  `json:"term_id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	Term_group int64  `json:"term_group"`
}

func (r *Wp_termsFields) RefColValue(name string) interface{} {
	switch name {
	case "term_id":
		return &r.Term_id

	case "name":
		return &r.Name

	case "slug":
		return &r.Slug

	case "term_group":
		return &r.Term_group

	default:
		return nil
	}
}

func (r *Wp_termsFields) ColValue(name string) interface{} {
	switch name {
	case "term_id":
		return r.Term_id

	case "name":
		return r.Name

	case "slug":
		return r.Slug

	case "term_group":
		return r.Term_group

	default:
		return nil
	}
}

func NewWp_terms(db *dbEngine.DB) (*Wp_terms, error) {
	table, ok := db.Tables["wp_terms"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "wp_terms"}
	}

	return &Wp_terms{
		Table: table,
	}, nil
}

func (t *Wp_terms) NewRecord() *Wp_termsFields {
	t.Record = &Wp_termsFields{}
	return t.Record
}

func (t *Wp_terms) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_terms) SelectSelfScanEach(ctx context.Context, each func(record *Wp_termsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_terms) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_terms) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
