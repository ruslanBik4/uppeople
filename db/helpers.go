package db

import (
	"strings"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
)

type SelectedUnit struct {
	Id    int32  `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type SelectedUnits []*SelectedUnit

func NewSelectedUnit(id int32, label string) *SelectedUnit {
	return &SelectedUnit{
		Id:    id,
		Label: label,
		Value: strings.ToLower(label),
	}
}

func (s *SelectedUnits) GetFields(columns []dbEngine.Column) []interface{} {
	p := &SelectedUnit{}
	r := make([]interface{}, 0)
	for _, col := range columns {
		switch col.Name() {
		case "id":
			r = append(r, &p.Id)
		case "label":
			r = append(r, &p.Label)
		case "value":
			r = append(r, &p.Value)
		default:
			logs.ErrorLog(dbEngine.ErrNotFoundColumn{Column: col.Name()}, "SelectedUnits. GetFields")
		}
	}

	*s = append(*s, p)

	return r
}
