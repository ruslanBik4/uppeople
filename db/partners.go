// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Partners struct {
	dbEngine.Table
	Record *PartnersFields
	rows   sql.Rows
}

type PartnersFields struct {
	Id         int64         `json:"id"`
	User_id    sql.NullInt64 `json:"user_id"`
	Vacancy_id sql.NullInt64 `json:"vacancy_id"`
}

func (r *PartnersFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "vacancy_id":
		return &r.Vacancy_id

	default:
		return nil
	}
}

func (r *PartnersFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "vacancy_id":
		return r.Vacancy_id

	default:
		return nil
	}
}

func NewPartners(db *dbEngine.DB) (*Partners, error) {
	table, ok := db.Tables[TABLE_PARTNERS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_PARTNERS}
	}

	return &Partners{
		Table: table,
	}, nil
}

func (t *Partners) NewRecord() *PartnersFields {
	t.Record = &PartnersFields{}
	return t.Record
}

func (t *Partners) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Partners) SelectSelfScanEach(ctx context.Context, each func(record *PartnersFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Partners) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Partners) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
