// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Candidates struct {
	dbEngine.Table
	Record *CandidatesFields
	rows   sql.Rows
}

type CandidatesFields struct {
	Id             int64          `json:"id"`
	Platform_id    sql.NullInt64  `json:"platform_id"`
	Name           sql.NullString `json:"name"`
	Salary         sql.NullString `json:"salary"`
	Email          sql.NullString `json:"email"`
	Mobile         sql.NullString `json:"mobile"`
	Skype          sql.NullString `json:"skype"`
	Link           sql.NullString `json:"link"`
	Linkedin       sql.NullString `json:"linkedin"`
	Str_companies  sql.NullString `json:"str_companies"`
	Status         sql.NullString `json:"status"`
	Tag_id         int64          `json:"tag_id"`
	Comments       sql.NullString `json:"comments"`
	Date           time.Time      `json:"date"`
	Recruter_id    sql.NullInt64  `json:"recruter_id"`
	Text_rezume    sql.NullString `json:"text_rezume"`
	Sfera          sql.NullString `json:"sfera"`
	Experience     sql.NullString `json:"experience"`
	Education      sql.NullString `json:"education"`
	Language       sql.NullString `json:"language"`
	Zapoln_profile sql.NullInt64  `json:"zapoln_profile"`
	File           sql.NullString `json:"file"`
	Avatar         sql.NullString `json:"avatar"`
	Seniority_id   sql.NullInt64  `json:"seniority_id"`
	Date_follow_up *time.Time     `json:"date_follow_up"`
}

func (r *CandidatesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "platform_id":
		return &r.Platform_id

	case "name":
		return &r.Name

	case "salary":
		return &r.Salary

	case "email":
		return &r.Email

	case "mobile":
		return &r.Mobile

	case "skype":
		return &r.Skype

	case "link":
		return &r.Link

	case "linkedin":
		return &r.Linkedin

	case "str_companies":
		return &r.Str_companies

	case "status":
		return &r.Status

	case "tag_id":
		return &r.Tag_id

	case "comments":
		return &r.Comments

	case "date":
		return &r.Date

	case "recruter_id":
		return &r.Recruter_id

	case "text_rezume":
		return &r.Text_rezume

	case "sfera":
		return &r.Sfera

	case "experience":
		return &r.Experience

	case "education":
		return &r.Education

	case "language":
		return &r.Language

	case "zapoln_profile":
		return &r.Zapoln_profile

	case "file":
		return &r.File

	case "avatar":
		return &r.Avatar

	case "seniority_id":
		return &r.Seniority_id

	case "date_follow_up":
		return &r.Date_follow_up

	default:
		return nil
	}
}

func (r *CandidatesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "platform_id":
		return r.Platform_id

	case "name":
		return r.Name

	case "salary":
		return r.Salary

	case "email":
		return r.Email

	case "mobile":
		return r.Mobile

	case "skype":
		return r.Skype

	case "link":
		return r.Link

	case "linkedin":
		return r.Linkedin

	case "str_companies":
		return r.Str_companies

	case "status":
		return r.Status

	case "tag_id":
		return r.Tag_id

	case "comments":
		return r.Comments

	case "date":
		return r.Date

	case "recruter_id":
		return r.Recruter_id

	case "text_rezume":
		return r.Text_rezume

	case "sfera":
		return r.Sfera

	case "experience":
		return r.Experience

	case "education":
		return r.Education

	case "language":
		return r.Language

	case "zapoln_profile":
		return r.Zapoln_profile

	case "file":
		return r.File

	case "avatar":
		return r.Avatar

	case "seniority_id":
		return r.Seniority_id

	case "date_follow_up":
		return r.Date_follow_up

	default:
		return nil
	}
}

func NewCandidates(db *dbEngine.DB) (*Candidates, error) {
	table, ok := db.Tables["candidates"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "candidates"}
	}

	return &Candidates{
		Table: table,
	}, nil
}

func (t *Candidates) NewRecord() *CandidatesFields {
	t.Record = &CandidatesFields{}
	return t.Record
}

func (t *Candidates) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Candidates) SelectSelfScanEach(ctx context.Context, each func(record *CandidatesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Candidates) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Candidates) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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