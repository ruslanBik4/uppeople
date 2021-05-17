// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Location_for_vacancies struct {
	dbEngine.Table
	Record *Location_for_vacanciesFields
	rows   sql.Rows
}

type Location_for_vacanciesFields struct {
	Id   int64          `json:"id"`
	Name sql.NullString `json:"name"`
}

func (r *Location_for_vacanciesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	default:
		return nil
	}
}

func (r *Location_for_vacanciesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	default:
		return nil
	}
}

func NewLocation_for_vacancies(db *dbEngine.DB) (*Location_for_vacancies, error) {
	table, ok := db.Tables[TABLE_LOCATION_FOR_VACANCIES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_LOCATION_FOR_VACANCIES}
	}

	return &Location_for_vacancies{
		Table: table,
	}, nil
}

func (t *Location_for_vacancies) NewRecord() *Location_for_vacanciesFields {
	t.Record = &Location_for_vacanciesFields{}
	return t.Record
}

func (t *Location_for_vacancies) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Location_for_vacancies) SelectSelfScanEach(ctx context.Context, each func(record *Location_for_vacanciesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Location_for_vacancies) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Location_for_vacancies) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
