// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_users struct {
	dbEngine.Table
	Record *Wp_usersFields
	rows   sql.Rows
}

type Wp_usersFields struct {
	Id                  int64     `json:"id"`
	User_login          string    `json:"user_login"`
	User_pass           string    `json:"user_pass"`
	User_nicename       string    `json:"user_nicename"`
	User_email          string    `json:"user_email"`
	User_url            string    `json:"user_url"`
	User_registered     time.Time `json:"user_registered"`
	User_activation_key string    `json:"user_activation_key"`
	User_status         int64     `json:"user_status"`
	Display_name        string    `json:"display_name"`
}

func (r *Wp_usersFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_login":
		return &r.User_login

	case "user_pass":
		return &r.User_pass

	case "user_nicename":
		return &r.User_nicename

	case "user_email":
		return &r.User_email

	case "user_url":
		return &r.User_url

	case "user_registered":
		return &r.User_registered

	case "user_activation_key":
		return &r.User_activation_key

	case "user_status":
		return &r.User_status

	case "display_name":
		return &r.Display_name

	default:
		return nil
	}
}

func (r *Wp_usersFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_login":
		return r.User_login

	case "user_pass":
		return r.User_pass

	case "user_nicename":
		return r.User_nicename

	case "user_email":
		return r.User_email

	case "user_url":
		return r.User_url

	case "user_registered":
		return r.User_registered

	case "user_activation_key":
		return r.User_activation_key

	case "user_status":
		return r.User_status

	case "display_name":
		return r.Display_name

	default:
		return nil
	}
}

func NewWp_users(db *dbEngine.DB) (*Wp_users, error) {
	table, ok := db.Tables[TABLE_WPUsers]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_WPUsers}
	}

	return &Wp_users{
		Table: table,
	}, nil
}

func (t *Wp_users) NewRecord() *Wp_usersFields {
	t.Record = &Wp_usersFields{}
	return t.Record
}

func (t *Wp_users) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_users) SelectSelfScanEach(ctx context.Context, each func(record *Wp_usersFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_users) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_users) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
