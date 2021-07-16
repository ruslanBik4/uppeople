package db

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"golang.org/x/net/context"
)

// TODO: add read table

type LogActions struct {
	dbEngine.Table
	Record *LogActionsFields
	rows   sql.Rows
}

type LogActionsFields struct {
	Id             int32  `json:"id"`
	Name           string `json:"name"`
	TextBeforeCand string `json:"text_before_cand"`
	ForCandidate   string `json:"for_candidate"`
	ForCompany     string `json:"for_company"`
	TextAfterCand  string `json:"text_after_cand"`
	IsInserText    bool   `json:"is_insert_text"`
}

type LogActionsIdMap map[string]LogActionsFields

func (r *LogActionsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "text_before_cand":
		return &r.TextBeforeCand

	case "for_candidate":
		return &r.ForCandidate

	case "for_company":
		return &r.ForCompany

	case "text_after_cand":
		return &r.TextAfterCand

	case "is_insert_text":
		return &r.IsInserText

	default:
		return nil
	}
}

func (r *LogActionsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "text_before_cand":
		return r.TextBeforeCand

	case "for_candidate":
		return r.ForCandidate

	case "for_company":
		return r.ForCompany

	case "text_after_cand":
		return r.TextAfterCand

	case "is_insert_text":
		return r.IsInserText

	default:
		return nil
	}
}

func NewLogActions(db *dbEngine.DB) (*LogActions, error) {
	table, ok := db.Tables[TABLE_LOG_ACTIONS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_LOG_ACTIONS}
	}

	return &LogActions{
		Table: table,
	}, nil
}

func (t *LogActions) NewRecord() *LogActionsFields {
	t.Record = &LogActionsFields{}
	return t.Record
}

func (t *LogActions) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *LogActions) SelectSelfScanEach(ctx context.Context, each func(record *LogActionsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *LogActions) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *LogActions) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetLogUpdateId() int32 {
	if logAction, ok := logActionsIds[LOG_UPDATE]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_UPDATE))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogInsertId() int32 {
	if logAction, ok := logActionsIds[LOG_INSERT]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_INSERT))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogPerformId() int32 {
	if logAction, ok := logActionsIds[LOG_PEFORM]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_PEFORM))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogDeleteId() int32 {
	if logAction, ok := logActionsIds[LOG_DELETE]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_DELETE))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogReContactId() int32 {
	if logAction, ok := logActionsIds[LOG_RE_CONTACT]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_RE_CONTACT))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogAddCommentId() int32 {
	if logAction, ok := logActionsIds[LOG_ADD_COMMENT]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_ADD_COMMENT))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogDelCommentId() int32 {
	if logAction, ok := logActionsIds[LOG_DEL_COMMENT]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_DEL_COMMENT))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogSendCVId() int32 {
	if logAction, ok := logActionsIds[LOG_SEND_CV]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_SEND_CV))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogAppointInterviewId() int32 {
	if logAction, ok := logActionsIds[LOG_APPOINT_INTERVIEW]; !ok {
		logs.ErrorLog(errors.Errorf("LogAction \"%s\" not found in database", LOG_APPOINT_INTERVIEW))
		return -1
	} else {
		return logAction.Id
	}
}

func GetLogActionFromId(id int32) *LogActionsFields {
	for _, logAction := range logActionsIds {
		if logAction.Id == id {
			return &logAction
		}
	}

	return nil
}

func GetLogActionsAsSelectedUnits() SelectedUnits {
	return logActionsSelected
}

func initLogActionsIds(ctx context.Context, db *dbEngine.DB) (err error) {
	logActionsIds = LogActionsIdMap{}
	logActionsTable, err := NewLogActions(db)
	if err != nil {
		logs.ErrorLog(err, "cannot get %s table", TABLE_LOG_ACTIONS)
		return err
	}

	err = logActionsTable.SelectSelfScanEach(ctx,
		func(record *LogActionsFields) error {
			logActionsIds[record.Name] = *record
			logActionsSelected = append(logActionsSelected, NewSelectedUnit(record.Id, record.Name))
			return nil
		},
		dbEngine.OrderBy("id"),
	)

	if err != nil {
		logs.ErrorLog(err, "while reading logActions from db to logActionsIds(db.LogActionsIdMap)")
	}

	return
}
