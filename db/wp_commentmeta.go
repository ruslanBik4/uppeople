// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_commentmeta struct {
	dbEngine.Table
	Record *Wp_commentmetaFields
	rows   sql.Rows
}

type Wp_commentmetaFields struct {
	Meta_id    int64          `json:"meta_id"`
	Comment_id float64        `json:"comment_id"`
	Meta_key   sql.NullString `json:"meta_key"`
	Meta_value sql.NullString `json:"meta_value"`
}

func (r *Wp_commentmetaFields) RefColValue(name string) interface{} {
	switch name {
	case "meta_id":
		return &r.Meta_id

	case "comment_id":
		return &r.Comment_id

	case "meta_key":
		return &r.Meta_key

	case "meta_value":
		return &r.Meta_value

	default:
		return nil
	}
}

func (r *Wp_commentmetaFields) ColValue(name string) interface{} {
	switch name {
	case "meta_id":
		return r.Meta_id

	case "comment_id":
		return r.Comment_id

	case "meta_key":
		return r.Meta_key

	case "meta_value":
		return r.Meta_value

	default:
		return nil
	}
}

func NewWp_commentmeta(db *dbEngine.DB) (*Wp_commentmeta, error) {
	table, ok := db.Tables[TABLE_WPCommentMeta]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_WPCommentMeta}
	}

	return &Wp_commentmeta{
		Table: table,
	}, nil
}

func (t *Wp_commentmeta) NewRecord() *Wp_commentmetaFields {
	t.Record = &Wp_commentmetaFields{}
	return t.Record
}

func (t *Wp_commentmeta) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_commentmeta) SelectSelfScanEach(ctx context.Context, each func(record *Wp_commentmetaFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_commentmeta) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_commentmeta) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
