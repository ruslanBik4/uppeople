package main

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
)

const TagsTable = "tags"

type TagIdMap map[string]TagStruct
type TagStruct struct {
	Id       int32  `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	ParentId int32  `json:"parent_id"`
	OrderNum int32  `json:"order_num"`
}

var TagsNames = map[string]string{
	"first contact":             "FirstContact",
	"interested":                "Interested",
	"reject":                    "Reject",
	"no answer":                 "NoAnswer",
	"closed to offers":          "ClosedToOffers",
	"low salary rate":           "LowSalary",
	"was contacted earlier":     "WasContactedEarlier",
	"does not like the project": "DoesNotLikeProject",
	"terms don’t fit":           "TermsDoNotFit",
	"remote only":               "RemoteOnly",
	"does not fit":              "DoesNotFit",
}

func (tagsStr *TagStruct) GetFields([]dbEngine.Column) []interface{} {
	return []interface{}{
		&tagsStr.Id,
		&tagsStr.Name,
		&tagsStr.Color,
		&tagsStr.ParentId,
		&tagsStr.OrderNum,
	}
}
