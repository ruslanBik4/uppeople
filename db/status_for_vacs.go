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
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_INTERVIEW]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_INTERVIEW))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdTest() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_TEST]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_TEST))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdFinalInterview() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_FINAL_INTERVIEW]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_FINAL_INTERVIEW))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdOffer() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_OFFER]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_OFFER))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdHired() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_HIRED]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_HIRED))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdWR() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_WR]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_WR))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdReview() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_REVIEW]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_REVIEW))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdRejected() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_REJECTED]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_REJECTED))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacIdOnHold() int32 {
	if statForVac, ok := statusesForVacIds[STATUS_FOR_VAC_ON_HOLD]; !ok {
		logs.ErrorLog(errors.Errorf("StatusForVac \"%s\" not found in database", STATUS_FOR_VAC_ON_HOLD))
		return -1
	} else {
		return statForVac.Id
	}
}

func GetStatusForVacFromId(id int32) *StatusForVacsFields {
	for _, statForVac := range statusesForVacIds {
		if statForVac.Id == id {
			return &statForVac
		}
	}

	return nil
}

func GetStatusForVacAsSelectedUnits() SelectedUnits {
	return statusesForVacIdsAsSU
}

func initStatusesForVacIds(ctx context.Context, db *dbEngine.DB) (err error) {
	statusesForVacIds = StatusForVacIdMap{}
	statusesForVacsTable, err := NewStatusForVacs(db)
	if err != nil {
		logs.ErrorLog(err, "cannot get %s table", TABLE_STATUS_FOR_VACS)
		return err
	}

	statusesForVacIdsAsSU = make(SelectedUnits, 0)
	err = statusesForVacsTable.SelectSelfScanEach(ctx,
		func(record *StatusForVacsFields) error {
			statusesForVacIds[record.Status] = *record
			statusesForVacIdsAsSU = append(statusesForVacIdsAsSU, &SelectedUnit{
				Id:    record.Id,
				Label: record.Status,
				Value: strings.ToLower(record.Status),
			})

			return nil
		},
		dbEngine.OrderBy("order_num"),
	)
	if err != nil {
		logs.ErrorLog(err, "while reading statuses for vacancies from db to statusesForVacIds(db.StatusForVacIdMap)")
	}

	return
}
