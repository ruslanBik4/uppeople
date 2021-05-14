// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Sended_emails struct {
	dbEngine.Table
	Record *Sended_emailsFields
	rows   sql.Rows
}

type Sended_emailsFields struct {
	Id          int64          `json:"id"`
	User_id     sql.NullInt64  `json:"user_id"`
	Company_id  sql.NullInt64  `json:"company_id"`
	Subject     sql.NullString `json:"subject"`
	Emails      sql.NullString `json:"emails"`
	Text_emails sql.NullString `json:"text_emails"`
	Date_create time.Time      `json:"date_create"`
	Meet_id     sql.NullInt64  `json:"meet_id"`
}

func (r *Sended_emailsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "user_id":
		return &r.User_id

	case "company_id":
		return &r.Company_id

	case "subject":
		return &r.Subject

	case "emails":
		return &r.Emails

	case "text_emails":
		return &r.Text_emails

	case "date_create":
		return &r.Date_create

	case "meet_id":
		return &r.Meet_id

	default:
		return nil
	}
}

func (r *Sended_emailsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "user_id":
		return r.User_id

	case "company_id":
		return r.Company_id

	case "subject":
		return r.Subject

	case "emails":
		return r.Emails

	case "text_emails":
		return r.Text_emails

	case "date_create":
		return r.Date_create

	case "meet_id":
		return r.Meet_id

	default:
		return nil
	}
}

func NewSended_emails(db *dbEngine.DB) (*Sended_emails, error) {
	table, ok := db.Tables[TABLE_SentEmails]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_SentEmails}
	}

	return &Sended_emails{
		Table: table,
	}, nil
}

func (t *Sended_emails) NewRecord() *Sended_emailsFields {
	t.Record = &Sended_emailsFields{}
	return t.Record
}

func (t *Sended_emails) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Sended_emails) SelectSelfScanEach(ctx context.Context, each func(record *Sended_emailsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Sended_emails) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Sended_emails) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
