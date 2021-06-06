// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Vacancies struct {
	dbEngine.Table
	Record *VacanciesFields
	rows   sql.Rows
}

type VacanciesFields struct {
	Id          int32          `json:"id"`
	CompanyId   int32          `json:"company_id"`
	PlatformId  int32          `json:"platform_id"`
	UserIds     []int32        `json:"user_ids"`
	Name        sql.NullString `json:"name"`
	Description string         `json:"description"`
	Details     string         `json:"details"`
	Link        string         `json:"link"`
	File        sql.NullString `json:"file"`
	DateCreate  time.Time      `json:"date_create"`
	Ord         sql.NullInt64  `json:"ord"`
	Status      int32          `json:"status"`
	SeniorityId int32          `json:"seniority_id"`
	Salary      int32          `json:"salary"`
	LocationId  int32          `json:"location_id"`
}

func (r *VacanciesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "company_id":
		return &r.CompanyId

	case "platform_id":
		return &r.PlatformId

	case "user_ids":
		return &r.UserIds

	case "name":
		return &r.Name

	case "description":
		return &r.Description

	case "details":
		return &r.Details

	case "link":
		return &r.Link

	case "file":
		return &r.File

	case "date_create":
		return &r.DateCreate

	case "ord":
		return &r.Ord

	case "status":
		return &r.Status

	case "seniority_id":
		return &r.SeniorityId

	case "salary":
		return &r.Salary

	case "location_id":
		return &r.LocationId

	default:
		return nil
	}
}

func (r *VacanciesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "company_id":
		return r.CompanyId

	case "platform_id":
		return r.PlatformId

	case "user_ids":
		return r.UserIds

	case "name":
		return r.Name

	case "description":
		return r.Description

	case "details":
		return r.Details

	case "link":
		return r.Link

	case "file":
		return r.File

	case "date_create":
		return r.DateCreate

	case "ord":
		return r.Ord

	case "status":
		return r.Status

	case "seniority_id":
		return r.SeniorityId

	case "salary":
		return r.Salary

	case "location_id":
		return r.LocationId

	default:
		return nil
	}
}

func NewVacancies(db *dbEngine.DB) (*Vacancies, error) {
	table, ok := db.Tables[TABLE_VACANCIES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_VACANCIES}
	}

	return &Vacancies{
		Table: table,
	}, nil
}

func (t *Vacancies) NewRecord() *VacanciesFields {
	t.Record = &VacanciesFields{}
	return t.Record
}

func (t *Vacancies) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Vacancies) SelectSelfScanEach(ctx context.Context, each func(record *VacanciesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Vacancies) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Vacancies) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
