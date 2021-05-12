// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"golang.org/x/net/context"
)

type Seniorities struct {
	dbEngine.Table
	Record *SenioritiesFields
	rows   sql.Rows
}

type SenioritiesFields struct {
	Id    int32          `json:"id"`
	Nazva sql.NullString `json:"nazva"`
}

type SeniorityIdMap map[string]SenioritiesFields

func (r *SenioritiesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "nazva":
		return &r.Nazva

	default:
		return nil
	}
}

func (r *SenioritiesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "nazva":
		return r.Nazva

	default:
		return nil
	}
}

func NewSeniorities(db *dbEngine.DB) (*Seniorities, error) {
	table, ok := db.Tables[TableSeniorities]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableSeniorities}
	}

	return &Seniorities{
		Table: table,
	}, nil
}

func (t *Seniorities) NewRecord() *SenioritiesFields {
	t.Record = &SenioritiesFields{}
	return t.Record
}

func (t *Seniorities) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Seniorities) SelectSelfScanEach(ctx context.Context, each func(record *SenioritiesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Seniorities) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Seniorities) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetSeniorityIdJun() int32 {
	if sen, ok := seniorityIds[SeniorityJun]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SeniorityJun))
	}

	return -1
}

func GetSeniorityIdMid() int32 {
	if sen, ok := seniorityIds[SeniorityMid]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SeniorityMid))
	}

	return -1
}

func GetSeniorityIdSen() int32 {
	if sen, ok := seniorityIds[SenioritySen]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SenioritySen))
	}

	return -1
}

func GetSeniorityIdLead() int32 {
	if sen, ok := seniorityIds[SeniorityLead]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SeniorityLead))
	}

	return -1
}

func GetSeniorityArchitect() int32 {
	if sen, ok := seniorityIds[SeniorityArchitect]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SeniorityArchitect))
	}

	return -1
}

func GetSeniorityIdJunMid() int32 {
	if sen, ok := seniorityIds[SeniorityJunMid]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SeniorityJunMid))
	}

	return -1
}

func GetSeniorityIdMidSen() int32 {
	if sen, ok := seniorityIds[SeniorityMidSen]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SeniorityMidSen))
	}

	return -1
}

func GetSeniorityIdSenLead() int32 {
	if sen, ok := seniorityIds[SenioritySenLead]; ok {
		return sen.Id
	} else {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SenioritySenLead))
	}

	return -1
}

func GetSeniorityFromId(id int32) *SenioritiesFields {
	for _, sen := range seniorityIds {
		if sen.Id == id {
			return &sen
		}
	}
	return nil
}
