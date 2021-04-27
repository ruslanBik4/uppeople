// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Companies struct {
	dbEngine.Table
	Record *CompaniesFields
	rows   sql.Rows
}

type CompaniesFields struct {
	Id              int64          `json:"id"`
	Name            sql.NullString `json:"name"`
	SendDetails     sql.NullString `json:"send_details"`
	InterviewDetail sql.NullString `json:"interview_detail"`
	Cooperation     sql.NullString `json:"cooperation"`
	Contact         sql.NullString `json:"contact"`
	About           sql.NullString `json:"about"`
	Map             sql.NullString `json:"map"`
	Phone           sql.NullString `json:"phone"`
	Email           sql.NullString `json:"email"`
	Skype           sql.NullString `json:"skype"`
	Logo            sql.NullString `json:"logo"`
	Address         sql.NullString `json:"address"`
	EmailTemplate   sql.NullString `json:"email_template"`
	ManagerId       sql.NullInt64  `json:"manager_id"`
}

func (r *CompaniesFields) GetValue() interface{} {
	return r
}

func (r *CompaniesFields) NewValue() interface{} {
	return &CompaniesFields{}
}

func (r *CompaniesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "send_details":
		return &r.SendDetails

	case "interview_detail":
		return &r.InterviewDetail

	case "cooperation":
		return &r.Cooperation

	case "contact":
		return &r.Contact

	case "about":
		return &r.About

	case "map":
		return &r.Map

	case "phone":
		return &r.Phone

	case "email":
		return &r.Email

	case "skype":
		return &r.Skype

	case "logo":
		return &r.Logo

	case "address":
		return &r.Address

	case "email_template":
		return &r.EmailTemplate

	case "manager_id":
		return &r.ManagerId

	default:
		return nil
	}
}

func (r *CompaniesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "send_details":
		return r.SendDetails

	case "interview_detail":
		return r.InterviewDetail

	case "cooperation":
		return r.Cooperation

	case "contact":
		return r.Contact

	case "about":
		return r.About

	case "map":
		return r.Map

	case "phone":
		return r.Phone

	case "email":
		return r.Email

	case "skype":
		return r.Skype

	case "logo":
		return r.Logo

	case "address":
		return r.Address

	case "email_template":
		return r.EmailTemplate

	case "manager_id":
		return r.ManagerId

	default:
		return nil
	}
}

func NewCompanies(db *dbEngine.DB) (*Companies, error) {
	table, ok := db.Tables[TableCompanies]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableCompanies}
	}

	return &Companies{
		Table: table,
	}, nil
}

func (t *Companies) NewRecord() *CompaniesFields {
	t.Record = &CompaniesFields{}
	return t.Record
}

func (t *Companies) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Companies) SelectSelfScanEach(ctx context.Context, each func(record *CompaniesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Companies) SelectSelfScanAll(ctx context.Context, each func(record *CompaniesFields) error,
	Options ...dbEngine.BuildSqlOptions) ([]*CompaniesFields, error) {
	rows := make([]*CompaniesFields, 0)
	err := t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}
			rows = append(rows, t.Record)

			return nil
		}, t, Options...)
	if err != nil {
		return nil, errors.Wrap(err, "SelectAndScanEach")
	}

	return rows, nil
}

func (t *Companies) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Companies) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
