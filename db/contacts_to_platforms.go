// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Contacts_to_platforms struct {
	dbEngine.Table
	Record *Contacts_to_platformsFields
	rows   sql.Rows
}

type Contacts_to_platformsFields struct {
	Id          int64         `json:"id"`
	Contact_id  sql.NullInt64 `json:"contact_id"`
	Platform_id sql.NullInt64 `json:"platform_id"`
}

func (r *Contacts_to_platformsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "contact_id":
		return &r.Contact_id

	case "platform_id":
		return &r.Platform_id

	default:
		return nil
	}
}

func (r *Contacts_to_platformsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "contact_id":
		return r.Contact_id

	case "platform_id":
		return r.Platform_id

	default:
		return nil
	}
}

func NewContacts_to_platforms(db *dbEngine.DB) (*Contacts_to_platforms, error) {
	table, ok := db.Tables["contacts_to_platforms"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "contacts_to_platforms"}
	}

	return &Contacts_to_platforms{
		Table: table,
	}, nil
}

func (t *Contacts_to_platforms) NewRecord() *Contacts_to_platformsFields {
	t.Record = &Contacts_to_platformsFields{}
	return t.Record
}

func (t *Contacts_to_platforms) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Contacts_to_platforms) SelectSelfScanEach(ctx context.Context, each func(record *Contacts_to_platformsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Contacts_to_platforms) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Contacts_to_platforms) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
