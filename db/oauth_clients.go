// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Oauth_clients struct {
	dbEngine.Table
	Record *Oauth_clientsFields
	rows   sql.Rows
}

type Oauth_clientsFields struct {
	Id                     int64         `json:"id"`
	User_id                sql.NullInt64 `json:"user_id"`
	Name                   string        `json:"name"`
	Secret                 string        `json:"secret"`
	Redirect               string        `json:"redirect"`
	Personal_access_client bool          `json:"personal_access_client"`
	Password_client        bool          `json:"password_client"`
	Revoked                bool          `json:"revoked"`
	Created_at             time.Time     `json:"created_at"`
	Updated_at             time.Time     `json:"updated_at"`
}

func (r *Oauth_clientsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "name":
		return &r.Name

	case "secret":
		return &r.Secret

	case "redirect":
		return &r.Redirect

	case "personal_access_client":
		return &r.Personal_access_client

	case "password_client":
		return &r.Password_client

	case "revoked":
		return &r.Revoked

	case "created_at":
		return &r.Created_at

	case "updated_at":
		return &r.Updated_at

	default:
		return nil
	}
}

func (r *Oauth_clientsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "name":
		return r.Name

	case "secret":
		return r.Secret

	case "redirect":
		return r.Redirect

	case "personal_access_client":
		return r.Personal_access_client

	case "password_client":
		return r.Password_client

	case "revoked":
		return r.Revoked

	case "created_at":
		return r.Created_at

	case "updated_at":
		return r.Updated_at

	default:
		return nil
	}
}

func NewOauth_clients(db *dbEngine.DB) (*Oauth_clients, error) {
	table, ok := db.Tables[TABLE_OauthClients]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_OauthClients}
	}

	return &Oauth_clients{
		Table: table,
	}, nil
}

func (t *Oauth_clients) NewRecord() *Oauth_clientsFields {
	t.Record = &Oauth_clientsFields{}
	return t.Record
}

func (t *Oauth_clients) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Oauth_clients) SelectSelfScanEach(ctx context.Context, each func(record *Oauth_clientsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Oauth_clients) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Oauth_clients) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
