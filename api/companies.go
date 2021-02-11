// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type SearchCompany struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	IsActive       string `json:"IsActive"`
	Skype          string `json:"skype"`
	Phone          string `json:"phone"`
	WithRecruiters string `json:"WithRecruiters"`
}

func (s *SearchCompany) GetValue() interface{} {
	return s
}

func (s *SearchCompany) NewValue() interface{} {
	return &SearchCompany{}
}

func HandleAllCompanies(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	offset := 0
	id, ok := ctx.UserValue(ParamPageNum.Name).(int)
	if ok && id > 1 {
		offset = id * pageItem
	}

	companies, _ := db.NewCompanies(DB)
	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy("name"),
		dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	}

	dto, ok := ctx.UserValue(apis.JSONParams).(*SearchCompany)
	if ok {
		args := make([]interface{}, 0)
		where := make([]string, 0)
		if dto.Name > "" {
			where = append(where, "~name")
			args = append(args, dto.Name)
		}
		if dto.Email > "" {
			where = append(where, "~email")
			args = append(args, dto.Name)
		}
		if dto.Skype > "" {
			where = append(where, "~skype")
			args = append(args, dto.Name)
		}
		if dto.Phone > "" {
			where = append(where, "~phone")
			args = append(args, dto.Name)
		}
		options = append(options,
			dbEngine.WhereForSelect(where...),
			dbEngine.ArgsForSelect(args...),
		)
	}

	return companies.SelectSelfScanAll(ctx, options...)
}
