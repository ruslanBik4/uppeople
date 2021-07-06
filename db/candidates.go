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
	Id           int32          `json:"id"`
	Platforms    []int32        `json:"platforms"`
	Name         string         `json:"name"`
	Salary       int32          `json:"salary"`
	Email        string         `json:"email"`
	Phone        string         `json:"phone"`
	Skype        string         `json:"skype"`
	Link         string         `json:"link"`
	Linkedin     sql.NullString `json:"linkedin"`
	Status       string         `json:"status"`
	Tag_id       int32          `json:"tag_id"`
	Comments     string         `json:"comment"`
	Date         time.Time      `json:"date"`
	RecruterId   int32          `json:"recruter_id"`
	Cv           string         `json:"cv"`
	Experience   string         `json:"experience"`
	Education    string         `json:"education"`
	IdLanguages  int32          `json:"id_languages"`
	File         string         `json:"file"`
	Avatar       []byte         `json:"avatar"`
	SeniorityId  int32          `json:"seniority_id"`
	DateFollowUp *time.Time     `json:"date_follow_up"`
	Vacancies    []int32        `json:"vacancies"`
}

func (r *CandidatesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "platforms":
		return &r.Platforms

	case "name":
		return &r.Name

	case "salary":
		return &r.Salary

	case "email":
		return &r.Email

	case "phone":
		return &r.Phone

	case "skype":
		return &r.Skype

	case "link":
		return &r.Link

	case "linkedin":
		return &r.Linkedin

	case "status":
		return &r.Status

	case "tag_id":
		return &r.Tag_id

	case "comments":
		return &r.Comments

	case "date":
		return &r.Date

	case "recruter_id":
		return &r.RecruterId

	case "text_rezume":
		return &r.Cv

	case "experience":
		return &r.Experience

	case "education":
		return &r.Education

	case "id_languages":
		return &r.IdLanguages

	case "file":
		return &r.File

	case "avatar":
		return &r.Avatar

	case "seniority_id":
		return &r.SeniorityId

	case "date_follow_up":
		return &r.DateFollowUp

	case "vacancies":
		return &r.Vacancies

	default:
		return nil
	}
}

func (r *CandidatesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "platforms":
		return r.Platforms

	case "name":
		return r.Name

	case "salary":
		return r.Salary

	case "email":
		return r.Email

	case "phone":
		return r.Phone

	case "skype":
		return r.Skype

	case "link":
		return r.Link

	case "linkedin":
		return r.Linkedin

	case "status":
		return r.Status

	case "tag_id":
		return r.Tag_id

	case "comments":
		return r.Comments

	case "date":
		return r.Date

	case "recruter_id":
		return r.RecruterId

	case "text_rezume":
		return r.Cv

	case "experience":
		return r.Experience

	case "education":
		return r.Education

	case "id_languages":
		return r.IdLanguages

	case "file":
		return r.File

	case "avatar":
		return r.Avatar

	case "seniority_id":
		return r.SeniorityId

	case "date_follow_up":
		return r.DateFollowUp

	case "vacancies":
		return r.Vacancies

	default:
		return nil
	}
}

func NewCandidates(db *dbEngine.DB) (*Candidates, error) {
	table, ok := db.Tables[TABLE_CANDIDATES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_CANDIDATES}
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
