// generate file
// don't edit
package db

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"golang.org/x/net/context"
	"strings"
)

type Tags struct {
	dbEngine.Table
	Record *TagsFields
	rows   sql.Rows
}

type TagsFields struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	ParentId int32  `json:"parent_id"`
	OrderNum int32  `json:"order_num"`
}

type TagIdMap map[string]TagsFields

func (r *TagsFields) GetFields(columns []dbEngine.Column) []interface{} {
	if len(columns) == 0 {
		return []interface{}{
			&r.Id,
			&r.Name,
			&r.Color,
			&r.ParentId,
			&r.OrderNum,
		}
	}

	v := make([]interface{}, len(columns))
	for i, col := range columns {
		v[i] = r.RefColValue(col.Name())
	}

	return v
}

func (r *TagsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "color":
		return &r.Color

	case "parent_id":
		return &r.ParentId

	case "order_num":
		return &r.OrderNum

	default:
		return nil
	}
}

func (r *TagsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "color":
		return r.Color

	case "parent_id":
		return r.ParentId

	case "order_num":
		return r.OrderNum

	default:
		return nil
	}
}

func NewTags(db *dbEngine.DB) (*Tags, error) {
	table, ok := db.Tables[TABLE_Tags]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_Tags}
	}

	return &Tags{
		Table: table,
	}, nil
}

func (t *Tags) NewRecord() *TagsFields {
	t.Record = &TagsFields{}
	return t.Record
}

func (t *Tags) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Tags) SelectSelfScanEach(ctx context.Context, each func(record *TagsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Tags) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Tags) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func GetTagIdFirstContact() int32 {
	if tag, ok := tagIds[TAG_FirstContact]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_FirstContact))
	}

	return -1
}

func GetTagIdInterested() int32 {
	if tag, ok := tagIds[TAG_Interested]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_Interested))
	}

	return -1
}

func GetTagIdReject() int32 {
	if tag, ok := tagIds[TAG_Reject]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_Reject))
	}

	return -1
}

func GetTagIdNoAnswer() int32 {
	if tag, ok := tagIds[TAG_NoAnswer]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_NoAnswer))
	}

	return -1
}

func GetTagIdClosedToOffers() int32 {
	if tag, ok := tagIds[TAG_ClosedToOffers]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_ClosedToOffers))
	}

	return -1
}

func GetTagIdLowSalary() int32 {
	if tag, ok := tagIds[TAG_LowSalary]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_LowSalary))
	}

	return -1
}

func GetTagIdWasContactedEarlier() int32 {
	if tag, ok := tagIds[TAG_WasContactedEarlier]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_WasContactedEarlier))
	}

	return -1
}

func GetTagIdDoesNotLikeProject() int32 {
	if tag, ok := tagIds[TAG_DoesNotLikeProject]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_DoesNotLikeProject))
	}

	return -1
}

func GetTagIdTermsDoNotFit() int32 {
	if tag, ok := tagIds[TAG_TermsDoNotFit]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_TermsDoNotFit))
	}

	return -1
}

func GetTagIdRemoteOnly() int32 {
	if tag, ok := tagIds[TAG_RemoteOnly]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_RemoteOnly))
	}

	return -1
}

func GetTagIdDoesNotFit() int32 {
	if tag, ok := tagIds[TAG_DoesNotFit]; ok {
		return tag.Id
	} else {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_DoesNotFit))
	}

	return -1
}

func GetTagFromId(id int32) *TagsFields {
	for _, tag := range tagIds {
		if tag.Id == id {
			return &tag
		}
	}

	return nil
}

func GetTagsAsSelectedUnits() (res SelectedUnits) {
	for _, tag := range tagIds {
		if tag.ParentId == 0 {
			res = append(res,
				&SelectedUnit{
					Id:    tag.Id,
					Label: tag.Name,
					Value: strings.ToLower(tag.Name),
				})
		}
	}

	if len(res) == 0 {
		return nil
	}

	return
}

func GetRejectReasonAsSelectedUnits() (res SelectedUnits) {
	for _, tag := range tagIds {
		if tag.ParentId == GetTagIdReject() {
			res = append(res,
				&SelectedUnit{
					Id:    tag.Id,
					Label: tag.Name,
					Value: strings.ToLower(tag.Name),
				})
		}
	}

	if len(res) == 0 {
		return nil
	}

	return
}
