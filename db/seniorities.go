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
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type SeniorityIdMap map[string]SenioritiesFields

func (r *SenioritiesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	default:
		return nil
	}
}

func (r *SenioritiesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	default:
		return nil
	}
}

func NewSeniorities(db *dbEngine.DB) (*Seniorities, error) {
	table, ok := db.Tables[TABLE_SENIORITIES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_SENIORITIES}
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
	if sen, ok := seniorityIds[SENIORITY_JUN]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_JUN))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityIdMid() int32 {
	if sen, ok := seniorityIds[SENIORITY_MID]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_MID))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityIdSen() int32 {
	if sen, ok := seniorityIds[SENIORITY_SEN]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_SEN))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityIdLead() int32 {
	if sen, ok := seniorityIds[SENIORITY_LEAD]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_LEAD))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityArchitect() int32 {
	if sen, ok := seniorityIds[SENIORITY_ARCHITECT]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_ARCHITECT))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityIdJunMid() int32 {
	if sen, ok := seniorityIds[SENIORITY_JUN_MID]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_JUN_MID))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityIdMidSen() int32 {
	if sen, ok := seniorityIds[SENIORITY_MID_SEN]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_MID_SEN))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityIdSenLead() int32 {
	if sen, ok := seniorityIds[SENIORITY_SEN_LEAD]; !ok {
		logs.ErrorLog(errors.Errorf("Seniority \"%s\" not found in database", SENIORITY_SEN_LEAD))
		return -1
	} else {
		return sen.Id
	}
}

func GetSeniorityFromId(id int32) *SenioritiesFields {
	for _, sen := range seniorityIds {
		if sen.Id == id {
			return &sen
		}
	}
	return nil
}

func GetSenioritiesAsSelectedUnits() SelectedUnits {
	return senioritySelected
}

func initSeniorityIds(ctx context.Context, db *dbEngine.DB) (err error) {
	seniorityIds = SeniorityIdMap{}
	seniorityTable, err := NewSeniorities(db)
	if err != nil {
		logs.ErrorLog(err, "cannot get %s table", TABLE_SENIORITIES)
		return err
	}

	err = seniorityTable.SelectSelfScanEach(ctx,
		func(record *SenioritiesFields) error {
			seniorityIds[record.Name] = *record
			senioritySelected = append(senioritySelected, NewSelectedUnit(record.Id, record.Name))
			return nil
		},
		dbEngine.OrderBy("id"),
	)

	if err != nil {
		logs.ErrorLog(err, "while reading seniorities from db to seniorityIds(db.SeniorityIdMap)")
	}

	return
}
