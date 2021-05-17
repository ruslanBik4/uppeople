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
	Id       int32  `json:"id"`
	Status   string `json:"status"`
	Color    string `json:"color"`
	OrderNum int64  `json:"order_num"`
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
	table, ok := db.Tables[TABLE_STATUS_FOR_VACS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_STATUS_FOR_VACS}
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
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_INTERVIEW]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_INTERVIEW))
	}

	return -1
}

func GetStatusForVacIdTest() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_TEST]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_TEST))
	}

	return -1
}

func GetStatusForVacIdFinalInterview() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_FINAL_INTERVIEW]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_FINAL_INTERVIEW))
	}

	return -1
}

func GetStatusForVacIdOffer() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_OFFER]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_OFFER))
	}

	return -1
}

func GetStatusForVacIdHired() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_HIRED]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_HIRED))
	}

	return -1
}

func GetStatusForVacIdWR() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_WR]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_WR))
	}

	return -1
}

func GetStatusForVacIdReview() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_REVIEW]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_REVIEW))
	}

	return -1
}

func GetStatusForVacIdRejected() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_REJECTED]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_REJECTED))
	}

	return -1
}

func GetStatusForVacIdOnHold() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_ON_HOLD]; ok {
		return statForVac.Id
	} else {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_ON_HOLD))
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
	if len(seniorityIds) == 0 {
		return nil
	}

	for _, statForVac := range statusesForVacIds {
		res = append(res,
			&SelectedUnit{
				Id:    statForVac.Id,
				Label: statForVac.Status,
				Value: strings.ToLower(statForVac.Status),
			})
	}
	return
}
