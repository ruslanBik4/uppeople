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

type Languages struct {
	dbEngine.Table
	Record *LanguagesFields
	rows   sql.Rows
}

type LanguagesFields struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
	Abbr string `json:"abbr"`
}

type LanguagesIdMap map[string]LanguagesFields

func (r *LanguagesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "abbr":
		return &r.Abbr

	default:
		return nil
	}
}

func (r *LanguagesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "abbr":
		return r.Abbr

	default:
		return nil
	}
}

func NewLanguages(db *dbEngine.DB) (*Languages, error) {
	table, ok := db.Tables[TABLE_LANGUAGES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_LANGUAGES}
	}

	return &Languages{
		Table: table,
	}, nil
}

func (t *Languages) NewRecord() *LanguagesFields {
	t.Record = &LanguagesFields{}
	return t.Record
}

func (t *Languages) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Languages) SelectSelfScanEach(ctx context.Context, each func(record *LanguagesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Languages) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Languages) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetBegginerId() int32 {
	if language, ok := languagesIds[LANG_BEGGINER]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_BEGGINER))
		return -1
	} else {
		return language.Id
	}
}

func GetElementaryId() int32 {
	if language, ok := languagesIds[LANG_ELEMENTARY]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_ELEMENTARY))
		return -1
	} else {
		return language.Id
	}
}

func GetPreIntermediateId() int32 {
	if language, ok := languagesIds[LANG_PRE_INT]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_PRE_INT))
		return -1
	} else {
		return language.Id
	}
}

func GetIntermediateId() int32 {
	if language, ok := languagesIds[LANG_INT]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_INT))
		return -1
	} else {
		return language.Id
	}
}
func GetUpperIntermediateId() int32 {
	if language, ok := languagesIds[LANG_UPPER_INT]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_UPPER_INT))
		return -1
	} else {
		return language.Id
	}
}

func GetAdvancedId() int32 {
	if language, ok := languagesIds[LANG_ADVANCED]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_ADVANCED))
		return -1
	} else {
		return language.Id
	}
}

func GetProficiencyId() int32 {
	if language, ok := languagesIds[LANG_PROF]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_PROF))
		return -1
	} else {
		return language.Id
	}
}

func GetUndefLanguagerId() int32 {
	if language, ok := languagesIds[LANG_UNDEF]; !ok {
		logs.ErrorLog(errors.Errorf("Language \"%s\" not found in database", LANG_UNDEF))
		return -1
	} else {
		return language.Id
	}
}

func GetLanguageFromId(id int32) *LanguagesFields {
	for _, language := range languagesIds {
		if language.Id == id {
			return &language
		}
	}

	return nil
}

func GetLanguagesAsSelectedUnits() SelectedUnits {
	return languagesSelected
}

func initLanguagesIds(ctx context.Context, db *dbEngine.DB) (err error) {
	languagesIds = LanguagesIdMap{}
	languagesTable, err := NewLanguages(db)
	if err != nil {
		logs.ErrorLog(err, "cannot get %s table", TABLE_LANGUAGES)
		return err
	}

	err = languagesTable.SelectSelfScanEach(ctx,
		func(record *LanguagesFields) error {
			languagesIds[record.Name] = *record
			languagesSelected = append(languagesSelected, NewSelectedUnit(record.Id, record.Name))
			return nil
		},
		dbEngine.OrderBy("id"),
	)

	if err != nil {
		logs.ErrorLog(err, "while reading languages from db to languagesIds(db.LanguagesIdMap)")
	}

	return
}
