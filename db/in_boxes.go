// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type In_boxes struct {
	dbEngine.Table
	Record *In_boxesFields
	rows   sql.Rows
}

type In_boxesFields struct {
	Id      int64          `json:"id"`
	Subject sql.NullString `json:"subject"`
	From    sql.NullString `json:"from"`
	Date    time.Time      `json:"date"`
	Email   sql.NullString `json:"email"`
}

func (r *In_boxesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "subject":
		return &r.Subject

	case "from":
		return &r.From

	case "date":
		return &r.Date

	case "email":
		return &r.Email

	default:
		return nil
	}
}

func (r *In_boxesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "subject":
		return r.Subject

	case "from":
		return r.From

	case "date":
		return r.Date

	case "email":
		return r.Email

	default:
		return nil
	}
}

func NewIn_boxes(db *dbEngine.DB) (*In_boxes, error) {
	table, ok := db.Tables[TABLE_InBoxes]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_InBoxes}
	}

	return &In_boxes{
		Table: table,
	}, nil
}

func (t *In_boxes) NewRecord() *In_boxesFields {
	t.Record = &In_boxesFields{}
	return t.Record
}

func (t *In_boxes) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *In_boxes) SelectSelfScanEach(ctx context.Context, each func(record *In_boxesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *In_boxes) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *In_boxes) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
