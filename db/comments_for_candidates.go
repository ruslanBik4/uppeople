// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Comments_for_candidates struct {
	dbEngine.Table
	Record *Comments_for_candidatesFields
	rows   sql.Rows
}

type Comments_for_candidatesFields struct {
	Id           int64          `json:"id"`
	User_id      sql.NullInt64  `json:"user_id"`
	Candidate_id sql.NullInt64  `json:"candidate_id"`
	Comments     sql.NullString `json:"comments"`
	Created_at   time.Time      `json:"created_at"`
	Updated_at   time.Time      `json:"updated_at"`
}

func (r *Comments_for_candidatesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "candidate_id":
		return &r.Candidate_id

	case "comments":
		return &r.Comments

	case "created_at":
		return &r.Created_at

	case "updated_at":
		return &r.Updated_at

	default:
		return nil
	}
}

func (r *Comments_for_candidatesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "candidate_id":
		return r.Candidate_id

	case "comments":
		return r.Comments

	case "created_at":
		return r.Created_at

	case "updated_at":
		return r.Updated_at

	default:
		return nil
	}
}

func NewComments_for_candidates(db *dbEngine.DB) (*Comments_for_candidates, error) {
	table, ok := db.Tables["comments_for_candidates"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "comments_for_candidates"}
	}

	return &Comments_for_candidates{
		Table: table,
	}, nil
}

func (t *Comments_for_candidates) NewRecord() *Comments_for_candidatesFields {
	t.Record = &Comments_for_candidatesFields{}
	return t.Record
}

func (t *Comments_for_candidates) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Comments_for_candidates) SelectSelfScanEach(ctx context.Context, each func(record *Comments_for_candidatesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Comments_for_candidates) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Comments_for_candidates) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
