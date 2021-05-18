// generate file
// don't edit
package db

import (
	"database/sql"
	"strings"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Platforms struct {
	dbEngine.Table
	Record *PlatformsFields
	rows   sql.Rows
}

type PlatformsFields struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

type PlatformsIdMap map[string]PlatformsFields

func (r *PlatformsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	default:
		return nil
	}
}

func (r *PlatformsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	default:
		return nil
	}
}

func NewPlatforms(db *dbEngine.DB) (*Platforms, error) {
	table, ok := db.Tables[TABLE_PLATFORMS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_PLATFORMS}
	}

	return &Platforms{
		Table: table,
	}, nil
}

func (t *Platforms) NewRecord() *PlatformsFields {
	t.Record = &PlatformsFields{}
	return t.Record
}

func (t *Platforms) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Platforms) SelectSelfScanEach(ctx context.Context, each func(record *PlatformsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Platforms) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Platforms) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetPlatformFromId(id int32) *PlatformsFields {
	for _, platform := range platformIds {
		if platform.Id == id {
			return &platform
		}
	}

	return nil
}

func GetPlatformsAsSelectedUnits() SelectedUnits {
	if len(platformIdsAsSU) > 0 {
		return platformIdsAsSU
	} else {
		if len(platformIds) == 0 {
			return nil
		}

		for _, platform := range platformIds {
			platformIdsAsSU = append(platformIdsAsSU,
				&SelectedUnit{
					Id:    platform.Id,
					Label: platform.Name,
					Value: strings.ToLower(platform.Name),
				})
		}

		if len(platformIdsAsSU) == 0 {
			return nil
		}
	}
	return platformIdsAsSU
}
