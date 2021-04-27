// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_usermeta struct {
	dbEngine.Table
	Record *Wp_usermetaFields
	rows   sql.Rows
}

type Wp_usermetaFields struct {
	Umeta_id   int64          `json:"umeta_id"`
	User_id    float64        `json:"user_id"`
	Meta_key   sql.NullString `json:"meta_key"`
	Meta_value sql.NullString `json:"meta_value"`
}

func (r *Wp_usermetaFields) RefColValue(name string) interface{} {
	switch name {
	case "umeta_id":
		return &r.Umeta_id

	case "user_id":
		return &r.User_id

	case "meta_key":
		return &r.Meta_key

	case "meta_value":
		return &r.Meta_value

	default:
		return nil
	}
}

func (r *Wp_usermetaFields) ColValue(name string) interface{} {
	switch name {
	case "umeta_id":
		return r.Umeta_id

	case "user_id":
		return r.User_id

	case "meta_key":
		return r.Meta_key

	case "meta_value":
		return r.Meta_value

	default:
		return nil
	}
}

func NewWp_usermeta(db *dbEngine.DB) (*Wp_usermeta, error) {
	table, ok := db.Tables[TableWPUserMeta]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableWPUserMeta}
	}

	return &Wp_usermeta{
		Table: table,
	}, nil
}

func (t *Wp_usermeta) NewRecord() *Wp_usermetaFields {
	t.Record = &Wp_usermetaFields{}
	return t.Record
}

func (t *Wp_usermeta) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_usermeta) SelectSelfScanEach(ctx context.Context, each func(record *Wp_usermetaFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_usermeta) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_usermeta) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
