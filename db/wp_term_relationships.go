// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_term_relationships struct {
	dbEngine.Table
	Record *Wp_term_relationshipsFields
	rows   sql.Rows
}

type Wp_term_relationshipsFields struct {
	Object_id        float64 `json:"object_id"`
	Term_taxonomy_id float64 `json:"term_taxonomy_id"`
	Term_order       int64   `json:"term_order"`
}

func (r *Wp_term_relationshipsFields) RefColValue(name string) interface{} {
	switch name {
	case "object_id":
		return &r.Object_id

	case "term_taxonomy_id":
		return &r.Term_taxonomy_id

	case "term_order":
		return &r.Term_order

	default:
		return nil
	}
}

func (r *Wp_term_relationshipsFields) ColValue(name string) interface{} {
	switch name {
	case "object_id":
		return r.Object_id

	case "term_taxonomy_id":
		return r.Term_taxonomy_id

	case "term_order":
		return r.Term_order

	default:
		return nil
	}
}

func NewWp_term_relationships(db *dbEngine.DB) (*Wp_term_relationships, error) {
	table, ok := db.Tables[TableWPTermRelationships]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableWPTermRelationships}
	}

	return &Wp_term_relationships{
		Table: table,
	}, nil
}

func (t *Wp_term_relationships) NewRecord() *Wp_term_relationshipsFields {
	t.Record = &Wp_term_relationshipsFields{}
	return t.Record
}

func (t *Wp_term_relationships) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_term_relationships) SelectSelfScanEach(ctx context.Context, each func(record *Wp_term_relationshipsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_term_relationships) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_term_relationships) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
