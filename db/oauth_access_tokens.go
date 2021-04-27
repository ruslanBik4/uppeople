// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Oauth_access_tokens struct {
	dbEngine.Table
	Record *Oauth_access_tokensFields
	rows   sql.Rows
}

type Oauth_access_tokensFields struct {
	Id         string         `json:"id"`
	User_id    sql.NullInt64  `json:"user_id"`
	Client_id  int64          `json:"client_id"`
	Name       sql.NullString `json:"name"`
	Scopes     sql.NullString `json:"scopes"`
	Revoked    bool           `json:"revoked"`
	Created_at time.Time      `json:"created_at"`
	Updated_at time.Time      `json:"updated_at"`
	Expires_at time.Time      `json:"expires_at"`
}

func (r *Oauth_access_tokensFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "client_id":
		return &r.Client_id

	case "name":
		return &r.Name

	case "scopes":
		return &r.Scopes

	case "revoked":
		return &r.Revoked

	case "created_at":
		return &r.Created_at

	case "updated_at":
		return &r.Updated_at

	case "expires_at":
		return &r.Expires_at

	default:
		return nil
	}
}

func (r *Oauth_access_tokensFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "client_id":
		return r.Client_id

	case "name":
		return r.Name

	case "scopes":
		return r.Scopes

	case "revoked":
		return r.Revoked

	case "created_at":
		return r.Created_at

	case "updated_at":
		return r.Updated_at

	case "expires_at":
		return r.Expires_at

	default:
		return nil
	}
}

func NewOauth_access_tokens(db *dbEngine.DB) (*Oauth_access_tokens, error) {
	table, ok := db.Tables[TableOauthAccessTokens]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableOauthAccessTokens}
	}

	return &Oauth_access_tokens{
		Table: table,
	}, nil
}

func (t *Oauth_access_tokens) NewRecord() *Oauth_access_tokensFields {
	t.Record = &Oauth_access_tokensFields{}
	return t.Record
}

func (t *Oauth_access_tokens) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Oauth_access_tokens) SelectSelfScanEach(ctx context.Context, each func(record *Oauth_access_tokensFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Oauth_access_tokens) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Oauth_access_tokens) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
