// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

func HandleAddCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*db.CompaniesFields)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)

	table, _ := db.NewCompanies(DB)
	for _, col := range table.Columns() {
		name := col.Name()
		if v := u.ColValue(name); v != nil {
			columns = append(columns, name)
			args = append(args, v)
		}
	}

	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCompany(ctx, DB, int32(i), "", CODE_LOG_INSERT)

	return createResult(i)
}

func HandleEditCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*db.CompaniesFields)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)

	table, _ := db.NewCompanies(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	for _, col := range table.Columns() {
		name := col.Name()
		if v := u.ColValue(name); v != table.Record.ColValue(name) {
			columns = append(columns, name)
			args = append(args, v)
		}
	}

	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(append(args, id)...),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCompany(ctx, DB, int32(i), "", CODE_LOG_UPDATE)

	return createResult(i)
}

func toLogCompany(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string, code int32) {
	user := auth.GetUserData(ctx)
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "company_id", "text", "date_create", "d_c",
			"kod_deystviya"),
		dbEngine.ArgsForSelect(user.Id, companyId,
			text,
			time.Now(),
			time.Now(),
			code))
}

func HandleCommentsCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	return DB.Conn.SelectToMaps(ctx,
		`select *, (select name from users u where u.id = user_id) as name
			 from comments_for_companies
			 where company_id = $1
			 order By time_create DESC`,
		id,
	)
}
