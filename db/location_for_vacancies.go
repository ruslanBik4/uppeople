// generate file
// don't edit
package db

import (
	"database/sql"
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
)

type LocationForVacancies struct {
	dbEngine.Table
	Record *LocationForVacanciesFields
	rows   sql.Rows
}

type LocationForVacanciesFields struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type LocationsIdMap map[string]LocationForVacanciesFields

func (r *LocationForVacanciesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	default:
		return nil
	}
}

func (r *LocationForVacanciesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	default:
		return nil
	}
}

func NewLocationForVacancies(db *dbEngine.DB) (*LocationForVacancies, error) {
	table, ok := db.Tables[TABLE_LOCATION_FOR_VACANCIES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_LOCATION_FOR_VACANCIES}
	}

	return &LocationForVacancies{
		Table: table,
	}, nil
}

func (t *LocationForVacancies) NewRecord() *LocationForVacanciesFields {
	t.Record = &LocationForVacanciesFields{}
	return t.Record
}

func (t *LocationForVacancies) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *LocationForVacancies) SelectSelfScanEach(ctx context.Context, each func(record *LocationForVacanciesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *LocationForVacancies) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *LocationForVacancies) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetRemoteId() int32 {
	if location, ok := locationsIds[LOCATION_REMOTE]; !ok {
		logs.ErrorLog(errors.Errorf("Location \"%s\" not found in database", LOCATION_REMOTE))
		return -1
	} else {
		return location.Id
	}
}

func GetKyivId() int32 {
	if location, ok := locationsIds[LOCATION_KYIV]; !ok {
		logs.ErrorLog(errors.Errorf("Location \"%s\" not found in database", LOCATION_KYIV))
		return -1
	} else {
		return location.Id
	}
}

func GetLvivId() int32 {
	if location, ok := locationsIds[LOCATION_LVIV]; !ok {
		logs.ErrorLog(errors.Errorf("Location \"%s\" not found in database", LOCATION_LVIV))
		return -1
	} else {
		return location.Id
	}
}

func GetOdessaId() int32 {
	if location, ok := locationsIds[LOCATION_ODESSA]; !ok {
		logs.ErrorLog(errors.Errorf("Location \"%s\" not found in database", LOCATION_ODESSA))
		return -1
	} else {
		return location.Id
	}
}

func GetKharkivId() int32 {
	if location, ok := locationsIds[LOCATION_KHARKIV]; !ok {
		logs.ErrorLog(errors.Errorf("Location \"%s\" not found in database", LOCATION_KHARKIV))
		return -1
	} else {
		return location.Id
	}
}

func GetDniproId() int32 {
	if location, ok := locationsIds[LOCATION_DNIPRO]; !ok {
		logs.ErrorLog(errors.Errorf("Location \"%s\" not found in database", LOCATION_DNIPRO))
		return -1
	} else {
		return location.Id
	}
}

func GetNopeLocationId() int32 {
	if location, ok := locationsIds[LOCATION_NOPE]; !ok {
		logs.ErrorLog(errors.Errorf("Location \"%s\" not found in database", LOCATION_NOPE))
		return -1
	} else {
		return location.Id
	}
}

func GetLocationFromId(id int32) *LocationForVacanciesFields {
	for _, location := range locationsIds {
		if location.Id == id {
			return &location
		}
	}

	return nil
}

func GetLocationsAsSelectedUnits() SelectedUnits {
	return locationsSelected
}

func initLocationsIds(ctx context.Context, db *dbEngine.DB) (err error) {
	locationsIds = LocationsIdMap{}
	locationsTable, err := NewLocationForVacancies(db)
	if err != nil {
		logs.ErrorLog(err, "cannot get %s table", TABLE_LOCATION_FOR_VACANCIES)
		return err
	}

	err = locationsTable.SelectSelfScanEach(ctx,
		func(record *LocationForVacanciesFields) error {
			locationsIds[record.Name] = *record
			locationsSelected = append(locationsSelected, NewSelectedUnit(record.Id, record.Name))
			return nil
		},
		dbEngine.OrderBy("id"),
	)

	if err != nil {
		logs.ErrorLog(err, "while reading locations from db to locationsIds(db.LocationsIdMap)")
	}

	return
}
