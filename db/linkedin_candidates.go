// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Linkedin_candidates struct {
	dbEngine.Table
	Record *Linkedin_candidatesFields
	rows   sql.Rows
}

type Linkedin_candidatesFields struct {
	Id          int64          `json:"id"`
	Platform_id sql.NullInt64  `json:"platform_id"`
	Name        sql.NullString `json:"name"`
	Email       sql.NullString `json:"email"`
	Mobile      sql.NullString `json:"mobile"`
	Skype       sql.NullString `json:"skype"`
	Linkedin    sql.NullString `json:"linkedin"`
	Status      sql.NullString `json:"status"`
	Recruter_id sql.NullInt64  `json:"recruter_id"`
	Text_rezume sql.NullString `json:"text_rezume"`
	Experience  sql.NullString `json:"experience"`
	Education   sql.NullString `json:"education"`
	Language    sql.NullString `json:"language"`
	Skills      sql.NullString `json:"skills"`
	Avatar      sql.NullString `json:"avatar"`
	Comment     sql.NullString `json:"comment"`
	Date_create time.Time      `json:"date_create"`
	Date_next   time.Time      `json:"date_next"`
}

func (r *Linkedin_candidatesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "platform_id":
		return &r.Platform_id

	case "name":
		return &r.Name

	case "email":
		return &r.Email

	case "mobile":
		return &r.Mobile

	case "skype":
		return &r.Skype

	case "linkedin":
		return &r.Linkedin

	case "status":
		return &r.Status

	case "recruter_id":
		return &r.Recruter_id

	case "text_rezume":
		return &r.Text_rezume

	case "experience":
		return &r.Experience

	case "education":
		return &r.Education

	case "language":
		return &r.Language

	case "skills":
		return &r.Skills

	case "avatar":
		return &r.Avatar

	case "comment":
		return &r.Comment

	case "date_create":
		return &r.Date_create

	case "date_next":
		return &r.Date_next

	default:
		return nil
	}
}

func (r *Linkedin_candidatesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "platform_id":
		return r.Platform_id

	case "name":
		return r.Name

	case "email":
		return r.Email

	case "mobile":
		return r.Mobile

	case "skype":
		return r.Skype

	case "linkedin":
		return r.Linkedin

	case "status":
		return r.Status

	case "recruter_id":
		return r.Recruter_id

	case "text_rezume":
		return r.Text_rezume

	case "experience":
		return r.Experience

	case "education":
		return r.Education

	case "language":
		return r.Language

	case "skills":
		return r.Skills

	case "avatar":
		return r.Avatar

	case "comment":
		return r.Comment

	case "date_create":
		return r.Date_create

	case "date_next":
		return r.Date_next

	default:
		return nil
	}
}

func NewLinkedin_candidates(db *dbEngine.DB) (*Linkedin_candidates, error) {
	table, ok := db.Tables["linkedin_candidates"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "linkedin_candidates"}
	}

	return &Linkedin_candidates{
		Table: table,
	}, nil
}

func (t *Linkedin_candidates) NewRecord() *Linkedin_candidatesFields {
	t.Record = &Linkedin_candidatesFields{}
	return t.Record
}

func (t *Linkedin_candidates) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Linkedin_candidates) SelectSelfScanEach(ctx context.Context, each func(record *Linkedin_candidatesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Linkedin_candidates) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Linkedin_candidates) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
