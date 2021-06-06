// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Faq_categories struct {
	dbEngine.Table
	Record *Faq_categoriesFields
	rows   sql.Rows
}

type Faq_categoriesFields struct {
	Id        int64          `json:"id"`
	Parent_id sql.NullInt64  `json:"parent_id"`
	Nazva     sql.NullString `json:"nazva"`
}

func (r *Faq_categoriesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "parent_id":
		return &r.Parent_id

	case "nazva":
		return &r.Nazva

	default:
		return nil
	}
}

func (r *Faq_categoriesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "parent_id":
		return r.Parent_id

	case "nazva":
		return r.Nazva

	default:
		return nil
	}
}

func NewFaq_categories(db *dbEngine.DB) (*Faq_categories, error) {
	table, ok := db.Tables[TABLE_FAQ_CATEGORIES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_FAQ_CATEGORIES}
	}

	return &Faq_categories{
		Table: table,
	}, nil
}

func (t *Faq_categories) NewRecord() *Faq_categoriesFields {
	t.Record = &Faq_categoriesFields{}
	return t.Record
}

func (t *Faq_categories) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Faq_categories) SelectSelfScanEach(ctx context.Context, each func(record *Faq_categoriesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Faq_categories) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Faq_categories) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
