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
	table, ok := db.Tables[TABLE_TAGS]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_TAGS}
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
	if tag, ok := tagIds[TAG_FIRST_CONTACT]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_FIRST_CONTACT))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdInterested() int32 {
	if tag, ok := tagIds[TAG_INTERESTED]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_INTERESTED))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdReject() int32 {
	if tag, ok := tagIds[TAG_REJECT]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_REJECT))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdNoAnswer() int32 {
	if tag, ok := tagIds[TAG_NO_ANSWER]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_NO_ANSWER))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdClosedToOffers() int32 {
	if tag, ok := tagIds[TAG_CLOSED_TO_OFFERS]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_CLOSED_TO_OFFERS))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdLowSalary() int32 {
	if tag, ok := tagIds[TAG_LOW_SALARY]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_LOW_SALARY))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdWasContactedEarlier() int32 {
	if tag, ok := tagIds[TAG_WAS_CONTACTED_EARLIER]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_WAS_CONTACTED_EARLIER))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdDoesNotLikeProject() int32 {
	if tag, ok := tagIds[TAG_DOES_NOT_LIKE_PROJECT]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_DOES_NOT_LIKE_PROJECT))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdTermsDoNotFit() int32 {
	if tag, ok := tagIds[TAG_TERMS_DO_NOT_FIT]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_TERMS_DO_NOT_FIT))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdRemoteOnly() int32 {
	if tag, ok := tagIds[TAG_REMOTE_ONLY]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_REMOTE_ONLY))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagIdDoesNotFit() int32 {
	if tag, ok := tagIds[TAG_DOES_NOT_FIT]; !ok {
		logs.ErrorLog(errors.Errorf("Tag \"%s\" not found in database", TAG_DOES_NOT_FIT))
		return -1
	} else {
		return tag.Id
	}
}

func GetTagFromId(id int32) *TagsFields {
	for _, tag := range tagIds {
		if tag.Id == id {
			return &tag
		}
	}

	return nil
}

func GetTagsAsSelectedUnits() SelectedUnits {
	if len(tagIdsAsSU) > 0 {
		return tagIdsAsSU
	}

	if len(tagIds) == 0 {
		return nil
	}

	for _, tag := range tagIds {
		if tag.ParentId == 0 {
			tagIdsAsSU = append(tagIdsAsSU,
				&SelectedUnit{
					Id:    tag.Id,
					Label: tag.Name,
					Value: strings.ToLower(tag.Name),
				})
		}
	}

	if len(tagIdsAsSU) == 0 {
		return nil
	}

	return tagIdsAsSU
}

func GetRejectReasonAsSelectedUnits() SelectedUnits {
	if len(reasonsIdsAsSU) > 0 {
		return reasonsIdsAsSU
	}

	if len(tagIds) == 0 {
		return nil
	}

	for _, tag := range tagIds {
		if tag.ParentId == GetTagIdReject() {
			reasonsIdsAsSU = append(reasonsIdsAsSU,
				&SelectedUnit{
					Id:    tag.Id,
					Label: tag.Name,
					Value: strings.ToLower(tag.Name),
				})
		}
	}

	if len(reasonsIdsAsSU) == 0 {
		return nil
	}

	return reasonsIdsAsSU
}
