// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_postmeta struct {
	dbEngine.Table
	Record *Wp_postmetaFields
	rows   sql.Rows
}

type Wp_postmetaFields struct {
	Meta_id    int64          `json:"meta_id"`
	Post_id    float64        `json:"post_id"`
	Meta_key   sql.NullString `json:"meta_key"`
	Meta_value sql.NullString `json:"meta_value"`
}

func (r *Wp_postmetaFields) RefColValue(name string) interface{} {
	switch name {
	case "meta_id":
		return &r.Meta_id

	case "post_id":
		return &r.Post_id

	case "meta_key":
		return &r.Meta_key

	case "meta_value":
		return &r.Meta_value

	default:
		return nil
	}
}

func (r *Wp_postmetaFields) ColValue(name string) interface{} {
	switch name {
	case "meta_id":
		return r.Meta_id

	case "post_id":
		return r.Post_id

	case "meta_key":
		return r.Meta_key

	case "meta_value":
		return r.Meta_value

	default:
		return nil
	}
}

func NewWp_postmeta(db *dbEngine.DB) (*Wp_postmeta, error) {
	table, ok := db.Tables[TABLE_WP_POST_META]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_WP_POST_META}
	}

	return &Wp_postmeta{
		Table: table,
	}, nil
}

func (t *Wp_postmeta) NewRecord() *Wp_postmetaFields {
	t.Record = &Wp_postmetaFields{}
	return t.Record
}

func (t *Wp_postmeta) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_postmeta) SelectSelfScanEach(ctx context.Context, each func(record *Wp_postmetaFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_postmeta) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_postmeta) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
