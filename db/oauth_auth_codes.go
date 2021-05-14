// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Oauth_auth_codes struct {
	dbEngine.Table
	Record *Oauth_auth_codesFields
	rows   sql.Rows
}

type Oauth_auth_codesFields struct {
	Id         string         `json:"id"`
	User_id    int64          `json:"user_id"`
	Client_id  int64          `json:"client_id"`
	Scopes     sql.NullString `json:"scopes"`
	Revoked    bool           `json:"revoked"`
	Expires_at time.Time      `json:"expires_at"`
}

func (r *Oauth_auth_codesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "client_id":
		return &r.Client_id

	case "scopes":
		return &r.Scopes

	case "revoked":
		return &r.Revoked

	case "expires_at":
		return &r.Expires_at

	default:
		return nil
	}
}

func (r *Oauth_auth_codesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "client_id":
		return r.Client_id

	case "scopes":
		return r.Scopes

	case "revoked":
		return r.Revoked

	case "expires_at":
		return r.Expires_at

	default:
		return nil
	}
}

func NewOauth_auth_codes(db *dbEngine.DB) (*Oauth_auth_codes, error) {
	table, ok := db.Tables[TABLE_OauthAuthCodes]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_OauthAuthCodes}
	}

	return &Oauth_auth_codes{
		Table: table,
	}, nil
}

func (t *Oauth_auth_codes) NewRecord() *Oauth_auth_codesFields {
	t.Record = &Oauth_auth_codesFields{}
	return t.Record
}

func (t *Oauth_auth_codes) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Oauth_auth_codes) SelectSelfScanEach(ctx context.Context, each func(record *Oauth_auth_codesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Oauth_auth_codes) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Oauth_auth_codes) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
