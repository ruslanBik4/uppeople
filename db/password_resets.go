// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Password_resets struct {
	dbEngine.Table
	Record *Password_resetsFields
	rows   sql.Rows
}

type Password_resetsFields struct {
	Email      string    `json:"email"`
	Token      string    `json:"token"`
	Created_at time.Time `json:"created_at"`
}

func (r *Password_resetsFields) RefColValue(name string) interface{} {
	switch name {
	case "email":
		return &r.Email

	case "token":
		return &r.Token

	case "created_at":
		return &r.Created_at

	default:
		return nil
	}
}

func (r *Password_resetsFields) ColValue(name string) interface{} {
	switch name {
	case "email":
		return r.Email

	case "token":
		return r.Token

	case "created_at":
		return r.Created_at

	default:
		return nil
	}
}

func NewPassword_resets(db *dbEngine.DB) (*Password_resets, error) {
	table, ok := db.Tables[TABLE_PASSWORD_RESETS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_PASSWORD_RESETS}
	}

	return &Password_resets{
		Table: table,
	}, nil
}

func (t *Password_resets) NewRecord() *Password_resetsFields {
	t.Record = &Password_resetsFields{}
	return t.Record
}

func (t *Password_resets) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Password_resets) SelectSelfScanEach(ctx context.Context, each func(record *Password_resetsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Password_resets) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Password_resets) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
