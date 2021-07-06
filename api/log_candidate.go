// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

func HandleReturnLogsForCand(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return DB.Conn.SelectToMaps(ctx, LOG_VIEW, ctx.UserValue("id"), true)
}

func HandleReturnLogsForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return DB.Conn.SelectToMaps(ctx, LOG_VIEW, ctx.UserValue("company_id"), false)
}

func toLogCandidate(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "text", "date_create",
			"action_code"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), candidateId,
			text,
			time.Now(),
			code))
}

func toLogCandidateVacancy(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId, companyId, vacancyId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "company_id", "vacancy_id", "text", "date_create",
			"action_code"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), candidateId, companyId, vacancyId,
			text,
			time.Now(),
			code))
}

func toLogCandidateStatus(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId, vacancyId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "vacancy_id", "text", "date_create",
			"action_code"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), candidateId, vacancyId,
			text,
			time.Now(),
			code))
}

func toLogCompany(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "company_id", "text", "date_create",
			"action_code"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), companyId,
			text,
			time.Now(),
			code))
}

func toLog(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, columns, args dbEngine.BuildSqlOptions) {
	log, _ := db.NewLogs(DB)
	_, err := log.Insert(ctx, columns, args)
	if err != nil {
		logs.ErrorLog(err, "toLog")
	}
}
