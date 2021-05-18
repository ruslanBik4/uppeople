// generate file
// don't edit
package db

import (
	"database/sql"
	"strings"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"golang.org/x/net/context"
)

type Statuses struct {
	dbEngine.Table
	Record *StatusesFields
	rows   sql.Rows
}

type StatusesFields struct {
	Id     int32  `json:"id"`
	Status string `json:"status"`
}

type StatusIdMap map[string]StatusesFields

func (r *StatusesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "status":
		return &r.Status

	default:
		return nil
	}
}

func (r *StatusesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "status":
		return r.Status

	default:
		return nil
	}
}

func NewStatuses(db *dbEngine.DB) (*Statuses, error) {
	table, ok := db.Tables[TABLE_STATUSES]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_STATUSES}
	}

	return &Statuses{
		Table: table,
	}, nil
}

func (t *Statuses) NewRecord() *StatusesFields {
	t.Record = &StatusesFields{}
	return t.Record
}

func (t *Statuses) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Statuses) SelectSelfScanEach(ctx context.Context, each func(record *StatusesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Statuses) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Statuses) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetStatusIdHot() int32 {
	if status, ok := statusesIds[STATUS_HOT]; ok {
		return status.Id
	} else {
		logs.ErrorLog(errors.Errorf("Status \"%s\" not found in database", STATUS_HOT))
	}

	return -1
}

func GetStatusIdOpen() int32 {
	if status, ok := statusesIds[STATUS_OPEN]; ok {
		return status.Id
	} else {
		logs.ErrorLog(errors.Errorf("Status \"%s\" not found in database", STATUS_OPEN))
	}

	return -1
}

func GetStatusIdClosed() int32 {
	if status, ok := statusesIds[STATUS_CLOSED]; ok {
		return status.Id
	} else {
		logs.ErrorLog(errors.Errorf("Status \"%s\" not found in database", STATUS_CLOSED))
	}

	return -1
}

func GetStatusIdPaused() int32 {
	if status, ok := statusesIds[STATUS_PAUSED]; ok {
		return status.Id
	} else {
		logs.ErrorLog(errors.Errorf("Status \"%s\" not found in database", STATUS_PAUSED))
	}

	return -1
}

func GetStatusFromId(id int32) *StatusesFields {
	for _, status := range statusesIds {
		if status.Id == id {
			return &status
		}
	}

	return nil
}

func GetStatusAsSelectedUnits() SelectedUnits {
	if len(statusesIdsAsSU) > 0 {
		return statusesIdsAsSU
	} else {
		if len(statusesIds) == 0 {
			return nil
		}

		for _, status := range statusesIds {
			statusesIdsAsSU = append(statusesIdsAsSU,
				&SelectedUnit{
					Id:    status.Id,
					Label: status.Status,
					Value: strings.ToLower(status.Status),
				})
		}

		if len(statusesIdsAsSU) == 0 {
			return nil
		}
	}

	return statusesIdsAsSU
}
