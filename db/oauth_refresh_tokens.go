// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Oauth_refresh_tokens struct {
	dbEngine.Table
	Record *Oauth_refresh_tokensFields
	rows   sql.Rows
}

type Oauth_refresh_tokensFields struct {
	Id              string    `json:"id"`
	Access_token_id string    `json:"access_token_id"`
	Revoked         bool      `json:"revoked"`
	Expires_at      time.Time `json:"expires_at"`
}

func (r *Oauth_refresh_tokensFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "access_token_id":
		return &r.Access_token_id

	case "revoked":
		return &r.Revoked

	case "expires_at":
		return &r.Expires_at

	default:
		return nil
	}
}

func (r *Oauth_refresh_tokensFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "access_token_id":
		return r.Access_token_id

	case "revoked":
		return r.Revoked

	case "expires_at":
		return r.Expires_at

	default:
		return nil
	}
}

func NewOauth_refresh_tokens(db *dbEngine.DB) (*Oauth_refresh_tokens, error) {
	table, ok := db.Tables[TABLE_OAUTH_REFRESH_TOKENS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_OAUTH_REFRESH_TOKENS}
	}

	return &Oauth_refresh_tokens{
		Table: table,
	}, nil
}

func (t *Oauth_refresh_tokens) NewRecord() *Oauth_refresh_tokensFields {
	t.Record = &Oauth_refresh_tokensFields{}
	return t.Record
}

func (t *Oauth_refresh_tokens) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Oauth_refresh_tokens) SelectSelfScanEach(ctx context.Context, each func(record *Oauth_refresh_tokensFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Oauth_refresh_tokens) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Oauth_refresh_tokens) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
