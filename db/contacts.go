// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Contacts struct {
	dbEngine.Table
	Record *ContactsFields
	rows   sql.Rows
}

type ContactsFields struct {
	Id         int64          `json:"id"`
	Company_id sql.NullInt64  `json:"company_id"`
	Name       sql.NullString `json:"name"`
	Email      sql.NullString `json:"email"`
	Phone      sql.NullString `json:"phone"`
	Skype      sql.NullString `json:"skype"`
	Default    sql.NullInt64  `json:"default"`
	Platforms  []int32        `json:"platforms"`
	NotVisible sql.NullInt64  `json:"not_visible"`
}

func (r *ContactsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "company_id":
		return &r.Company_id

	case "name":
		return &r.Name

	case "email":
		return &r.Email

	case "phone":
		return &r.Phone

	case "skype":
		return &r.Skype

	case "default_contact":
		return &r.Default

	case "platforms":
		return &r.Platforms

	case "not_visible":
		return &r.NotVisible

	default:
		return nil
	}
}

func (r *ContactsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "company_id":
		return r.Company_id

	case "name":
		return r.Name

	case "email":
		return r.Email

	case "phone":
		return r.Phone

	case "skype":
		return r.Skype

	case "default_contact":
		return r.Default

	case "platforms":
		return r.Platforms

	case "not_visible":
		return r.NotVisible

	default:
		return nil
	}
}

func NewContacts(db *dbEngine.DB) (*Contacts, error) {
	table, ok := db.Tables[TABLE_CONTACTS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_CONTACTS}
	}

	return &Contacts{
		Table: table,
	}, nil
}

func (t *Contacts) NewRecord() *ContactsFields {
	t.Record = &ContactsFields{}
	return t.Record
}

func (t *Contacts) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Contacts) SelectSelfScanEach(ctx context.Context, each func(record *ContactsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Contacts) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Contacts) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
