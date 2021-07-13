// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
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

func toLogCandidateInsert(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string) {
	toLogCandidate(ctx, DB, candidateId, text, db.GetLogInsertId())
}

func toLogCandidateUpdate(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text map[string]interface{}) {
	toLogCandidate(ctx, DB, candidateId, text, db.GetLogUpdateId())
}

func toLogCandidatePerform(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string) {
	toLogCandidate(ctx, DB, candidateId, text, db.GetLogPerformId())
}

func toLogCandidateDelete(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string) {
	toLogCandidate(ctx, DB, candidateId, text, db.GetLogDeleteId())
}

func toLogCandidateRecontact(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string) {
	toLogCandidate(ctx, DB, candidateId, text, db.GetLogReContactId())
}

func toLogCandidateAddComment(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string) {
	toLogCandidate(ctx, DB, candidateId, text, db.GetLogAddCommentId())
}

func toLogCandidateDelComment(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string) {
	toLogCandidate(ctx, DB, candidateId, text, db.GetLogDelCommentId())
}

func toLogCandidateSendCV(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId, companyId, vacancyId int32, text string) {
	toLogCandidateVacancy(ctx, DB, candidateId, companyId, vacancyId, text, db.GetLogSendCVId())
}

func toLogCandidateAppointInterview(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId, companyId, vacancyId int32, text string) {
	toLogCandidateVacancy(ctx, DB, candidateId, companyId, vacancyId, text, db.GetLogAppointInterviewId())
}

func toLogCompanyInsert(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string) {
	toLogCompany(ctx, DB, companyId, text, db.GetLogInsertId())
}

func toLogCompanyUpdate(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string) {
	toLogCompany(ctx, DB, companyId, text, db.GetLogUpdateId())
}

func toLogCompanyAddComment(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string) {
	toLogCompany(ctx, DB, companyId, text, db.GetLogAddCommentId())
}

func toLogCompanyDelComment(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string) {
	toLogCompany(ctx, DB, companyId, text, db.GetLogDelCommentId())
}

func toLogCompanyDelete(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string) {
	toLogCompany(ctx, DB, companyId, text, db.GetLogDeleteId())
}

func toLogCompanyPerform(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string) {
	toLogCompany(ctx, DB, companyId, text, db.GetLogPerformId())
}

func toLogVacancyInsert(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, DB, companyId, vacancyId, text, db.GetLogInsertId())
}

func toLogVacancyUpdate(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, DB, companyId, vacancyId, text, db.GetLogUpdateId())
}

func toLogVacancyPerform(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, DB, companyId, vacancyId, text, db.GetLogPerformId())
}

func toLogVacancyDelete(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, DB, companyId, vacancyId, text, db.GetLogDeleteId())
}

func toLogCandidateUpdateStatus(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId, vacancyId int32, text string) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "vacancy_id", "text", "date_create",
			"action_code"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), candidateId, vacancyId,
			text,
			time.Now(),
			db.GetLogUpdateId()))
}

func toLogCandidate(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text interface{}, code int32) {
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

func toLogCompany(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "company_id", "text", "date_create",
			"action_code"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), companyId,
			text,
			time.Now(),
			code))
}

func toLogVacancy(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId, vacancyId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "company_id", "vacancy_id", "text", "date_create",
			"action_code"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), companyId, vacancyId,
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

func loLogUpdateValues(columns []string, args []interface{}) string {
	if len(columns) > 0 {
		text := "{"
		for i, col := range columns {
			if i > 0 && i < len(columns)-1 {
				text += ", "
			}
			switch args[i].(type) {
			case int32:
				text += fmt.Sprintf("\""+"%s\":%v", col, args[i])
			case string:
				text += fmt.Sprintf("\""+"%s\":\"%v\"", col, args[i])
			case time.Time:
				text += fmt.Sprintf("\""+"%s\":\"%v\"", col, args[i])

			default:
				text += fmt.Sprintf("\""+"%s\":\"%v\"", col, args[i])
			}
		}
		text += "}"
		return text
	}

	return ""
}
