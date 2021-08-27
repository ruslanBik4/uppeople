// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Users struct {
	dbEngine.Table
	Record *UsersFields
	rows   sql.Rows
}

type UsersFields struct {
	Id          int32          `json:"id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	Isdel       bool           `json:"isdel"`
	RoleId      int32          `json:"role_id"`
	LastLogin   time.Time      `json:"last_login"`
	Schema      string         `json:"schema"`
	LastIP      string         `json:"last_ip"`
	Hash        int64          `json:"hash"`
	LastPage    sql.NullString `json:"last_page"`
	Address     string         `json:"address"`
	EmailPool   []string       `json:"emailpool"`
	Phone       sql.NullString `json:"phone"`
	Languages   []string       `json:"languages"`
	IdHomepages int32          `json:"id_homepages"`
	CreateAt    time.Time      `json:"createat"`
}

func (r *UsersFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "email":
		return &r.Email

	case "isdel":
		return &r.Isdel

	case "role_id":
		return &r.RoleId

	case "last_login":
		return &r.LastLogin

	case "hash":
		return &r.Hash

	case "last_page":
		return &r.LastPage

	case "schema":
		return &r.Schema

	case "last_ip":
		return &r.LastIP

	case "address":
		return &r.Address

	case "emailpool":
		return &r.EmailPool

	case "phone":
		return &r.Phone

	case "languages":
		return &r.Languages

	case "id_homepages":
		return &r.IdHomepages

	case "createat":
		return &r.CreateAt

	default:
		return nil
	}
}

func (r *UsersFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "email":
		return r.Email

	case "isdel":
		return r.Isdel

	case "role_id":
		return r.RoleId

	case "last_login":
		return r.LastLogin

	case "hash":
		return r.Hash

	case "last_page":
		return r.LastPage

	case "schema":
		return r.Schema

	case "last_ip":
		return r.LastIP

	case "address":
		return r.Address

	case "emailpool":
		return r.EmailPool

	case "phone":
		return r.Phone

	case "languages":
		return r.Languages

	case "id_homepages":
		return r.IdHomepages

	case "createat":
		return r.CreateAt

	default:
		return nil
	}
}

func NewUsers(db *dbEngine.DB) (*Users, error) {
	table, ok := db.Tables[TABLE_USERS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_USERS}
	}

	return &Users{
		Table: table,
	}, nil
}

func (t *Users) NewRecord() *UsersFields {
	t.Record = &UsersFields{}
	return t.Record
}

func (t *Users) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Users) SelectSelfScanEach(ctx context.Context, each func(record *UsersFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Users) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Users) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
