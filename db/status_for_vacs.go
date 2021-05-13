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

type StatusForVacs struct {
	dbEngine.Table
	Record *StatusForVacsFields
	rows   sql.Rows
}

type StatusForVacsFields struct {
	Id       int32          `json:"id"`
	Status   sql.NullString `json:"status"`
	Color    string         `json:"color"`
	OrderNum int64          `json:"order_num"`
}

type StatusForVacIdMap map[string]StatusForVacsFields

func (r *StatusForVacsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "status":
		return &r.Status

	case "color":
		return &r.Color

	case "order_num":
		return &r.OrderNum

	default:
		return nil
	}
}

func (r *StatusForVacsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "status":
		return r.Status

	case "color":
		return r.Color

	case "order_num":
		return r.OrderNum

	default:
		return nil
	}
}

func NewStatusForVacs(db *dbEngine.DB) (*StatusForVacs, error) {
	table, ok := db.Tables[TableStatusForVacs]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableStatusForVacs}
	}

	return &StatusForVacs{
		Table: table,
	}, nil
}

func (t *StatusForVacs) NewRecord() *StatusForVacsFields {
	t.Record = &StatusForVacsFields{}
	return t.Record
}

func (t *StatusForVacs) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *StatusForVacs) SelectSelfScanEach(ctx context.Context, each func(record *StatusForVacsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *StatusForVacs) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *StatusForVacs) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetStatusForVacIdInterview() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacInterview]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacInterview))
	}

	return -1
}

func GetStatusForVacIdTest() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacTest]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacTest))
	}

	return -1
}

func GetStatusForVacIdFinalInterview() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacFinalInterview]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacFinalInterview))
	}

	return -1
}

func GetStatusForVacIdOffer() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacOffer]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacOffer))
	}

	return -1
}

func GetStatusForVacIdHired() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacHired]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacHired))
	}

	return -1
}

func GetStatusForVacIdWR() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacWR]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacWR))
	}

	return -1
}

func GetStatusForVacIdReview() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacReview]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacReview))
	}

	return -1
}

func GetStatusForVacIdRejected() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacRejected]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacRejected))
	}

	return -1
}

func GetStatusForVacIdOnHold() int32 {
	if statForVac, ok := statusesForVacIds[StatusForVacOnHold]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", StatusForVacOnHold))
	}

	return -1
}

func GetStatusForVacFromId(id int32) *StatusForVacsFields {
	for _, statForVac := range statusesForVacIds {
		if statForVac.Id == id {
			return &statForVac
		}
	}

	return nil
}

func GetStatusForVacAsSelectedUnits() (res SelectedUnits) {
	for _, statForVac := range statusesForVacIds {
		res = append(res,
			&SelectedUnit{
				Id:    statForVac.Id,
				Label: statForVac.Status.String,
				Value: strings.ToLower(statForVac.Status.String),
			})
	}
	return
}
