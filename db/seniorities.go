// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Seniorities struct {
	dbEngine.Table
	Record *SenioritiesFields
	rows   sql.Rows
}

type SenioritiesFields struct {
	Id    int64          `json:"id"`
	Nazva sql.NullString `json:"nazva"`
}

func (r *SenioritiesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "nazva":
		return &r.Nazva

	default:
		return nil
	}
}

func (r *SenioritiesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "nazva":
		return r.Nazva

	default:
		return nil
	}
}

func NewSeniorities(db *dbEngine.DB) (*Seniorities, error) {
	table, ok := db.Tables[TableSeniorities]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableSeniorities}
	}

	return &Seniorities{
		Table: table,
	}, nil
}

func (t *Seniorities) NewRecord() *SenioritiesFields {
	t.Record = &SenioritiesFields{}
	return t.Record
}

func (t *Seniorities) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Seniorities) SelectSelfScanEach(ctx context.Context, each func(record *SenioritiesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Seniorities) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Seniorities) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
