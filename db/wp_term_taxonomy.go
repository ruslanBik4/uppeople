// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_term_taxonomy struct {
	dbEngine.Table
	Record *Wp_term_taxonomyFields
	rows   sql.Rows
}

type Wp_term_taxonomyFields struct {
	Term_taxonomy_id int64   `json:"term_taxonomy_id"`
	Term_id          float64 `json:"term_id"`
	Taxonomy         string  `json:"taxonomy"`
	Description      string  `json:"description"`
	Parent           float64 `json:"parent"`
	Count            int64   `json:"count"`
}

func (r *Wp_term_taxonomyFields) RefColValue(name string) interface{} {
	switch name {
	case "term_taxonomy_id":
		return &r.Term_taxonomy_id

	case "term_id":
		return &r.Term_id

	case "taxonomy":
		return &r.Taxonomy

	case "description":
		return &r.Description

	case "parent":
		return &r.Parent

	case "count":
		return &r.Count

	default:
		return nil
	}
}

func (r *Wp_term_taxonomyFields) ColValue(name string) interface{} {
	switch name {
	case "term_taxonomy_id":
		return r.Term_taxonomy_id

	case "term_id":
		return r.Term_id

	case "taxonomy":
		return r.Taxonomy

	case "description":
		return r.Description

	case "parent":
		return r.Parent

	case "count":
		return r.Count

	default:
		return nil
	}
}

func NewWp_term_taxonomy(db *dbEngine.DB) (*Wp_term_taxonomy, error) {
	table, ok := db.Tables[TableWPTermTaxonomy]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableWPTermTaxonomy}
	}

	return &Wp_term_taxonomy{
		Table: table,
	}, nil
}

func (t *Wp_term_taxonomy) NewRecord() *Wp_term_taxonomyFields {
	t.Record = &Wp_term_taxonomyFields{}
	return t.Record
}

func (t *Wp_term_taxonomy) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_term_taxonomy) SelectSelfScanEach(ctx context.Context, each func(record *Wp_term_taxonomyFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_term_taxonomy) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_term_taxonomy) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
