// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Favorites_platforms struct {
	dbEngine.Table
	Record *Favorites_platformsFields
	rows   sql.Rows
}

type Favorites_platformsFields struct {
	Id         int64         `json:"id"`
	Vacancy_id sql.NullInt64 `json:"vacancy_id"`
	User_id    sql.NullInt64 `json:"user_id"`
}

func (r *Favorites_platformsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "vacancy_id":
		return &r.Vacancy_id

	case "user_id":
		return &r.User_id

	default:
		return nil
	}
}

func (r *Favorites_platformsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "vacancy_id":
		return r.Vacancy_id

	case "user_id":
		return r.User_id

	default:
		return nil
	}
}

func NewFavorites_platforms(db *dbEngine.DB) (*Favorites_platforms, error) {
	table, ok := db.Tables[TABLE_FAVORITE_PLATFORMS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_FAVORITE_PLATFORMS}
	}

	return &Favorites_platforms{
		Table: table,
	}, nil
}

func (t *Favorites_platforms) NewRecord() *Favorites_platformsFields {
	t.Record = &Favorites_platformsFields{}
	return t.Record
}

func (t *Favorites_platforms) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Favorites_platforms) SelectSelfScanEach(ctx context.Context, each func(record *Favorites_platformsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Favorites_platforms) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Favorites_platforms) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
