// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Comments_for_companies struct {
	dbEngine.Table
	Record *Comments_for_companiesFields
	rows   sql.Rows
}

type Comments_for_companiesFields struct {
	Id          int64          `json:"id"`
	User_id     sql.NullInt64  `json:"user_id"`
	Company_id  sql.NullInt64  `json:"company_id"`
	Comments    sql.NullString `json:"comments"`
	Time_create sql.NullString `json:"time_create"`
}

func (r *Comments_for_companiesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "company_id":
		return &r.Company_id

	case "comments":
		return &r.Comments

	case "time_create":
		return &r.Time_create

	default:
		return nil
	}
}

func (r *Comments_for_companiesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "company_id":
		return r.Company_id

	case "comments":
		return r.Comments

	case "time_create":
		return r.Time_create

	default:
		return nil
	}
}

func NewComments_for_companies(db *dbEngine.DB) (*Comments_for_companies, error) {
	table, ok := db.Tables[TABLE_CommentsForCompanies]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_CommentsForCompanies}
	}

	return &Comments_for_companies{
		Table: table,
	}, nil
}

func (t *Comments_for_companies) NewRecord() *Comments_for_companiesFields {
	t.Record = &Comments_for_companiesFields{}
	return t.Record
}

func (t *Comments_for_companies) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Comments_for_companies) SelectSelfScanEach(ctx context.Context, each func(record *Comments_for_companiesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Comments_for_companies) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Comments_for_companies) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
