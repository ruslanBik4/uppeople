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

var (
	columnsForCandidateLog        = []string{"user_id", "candidate_id", "text", "action_code", "date_create"}
	columnsForCandidateVacancyLog = []string{"user_id", "candidate_id", "company_id", "vacancy_id", "text", "action_code", "date_create"}
	columnsForVacancyLog          = []string{"user_id", "company_id", "vacancy_id", "text", "action_code", "date_create"}
	columnsForCompanyLog          = []string{"user_id", "company_id", "text", "action_code", "date_create"}
	columnsForCandidateStatusLog  = []string{"user_id", "candidate_id", "vacancy_id", "text", "action_code", "date_create"}
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

func toLogCandidateInsert(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, text, db.GetLogInsertId())
}

func toLogCandidateUpdate(ctx *fasthttp.RequestCtx, candidateId int32, text map[string]interface{}) {
	toLogCandidate(ctx, candidateId, text, db.GetLogUpdateId())
}

func toLogCandidatePerform(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, text, db.GetLogPerformId())
}

func toLogCandidateDelete(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, text, db.GetLogDeleteId())
}

func toLogCandidateRecontact(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, text, db.GetLogReContactId())
}

func toLogCandidateAddComment(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, text, db.GetLogAddCommentId())
}

func toLogCandidateDelComment(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, text, db.GetLogDelCommentId())
}

func toLogCandidateSendCV(ctx *fasthttp.RequestCtx, candidateId, companyId, vacancyId int32, text string) {
	toLogCandidateVacancy(ctx, candidateId, companyId, vacancyId, text, db.GetLogSendCVId())
}

func toLogCandidateAppointInterview(ctx *fasthttp.RequestCtx, candidateId, companyId, vacancyId int32, text string) {
	toLogCandidateVacancy(ctx, candidateId, companyId, vacancyId, text, db.GetLogAppointInterviewId())
}

func toLogCompanyInsert(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, text, db.GetLogInsertId())
}

func toLogCompanyUpdate(ctx *fasthttp.RequestCtx, companyId int32, text map[string]interface{}) {
	toLogCompany(ctx, companyId, text, db.GetLogUpdateId())
}

func toLogCompanyAddComment(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, text, db.GetLogAddCommentId())
}

func toLogCompanyDelComment(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, text, db.GetLogDelCommentId())
}

func toLogCompanyDelete(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, text, db.GetLogDeleteId())
}

func toLogCompanyPerform(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, text, db.GetLogPerformId())
}

func toLogVacancyInsert(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, companyId, vacancyId, text, db.GetLogInsertId())
}

func toLogVacancyUpdate(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text map[string]interface{}) {
	toLogVacancy(ctx, companyId, vacancyId, text, db.GetLogUpdateId())
}

func toLogVacancyPerform(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, companyId, vacancyId, text, db.GetLogPerformId())
}

func toLogVacancyDelete(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, companyId, vacancyId, text, db.GetLogDeleteId())
}

func toLogCandidateUpdateStatus(ctx *fasthttp.RequestCtx, args ...interface{}) {
	args = append(args, db.GetLogUpdateId())
	toLog(ctx, columnsForCandidateStatusLog, args)
}

func toLogCandidate(ctx *fasthttp.RequestCtx, args ...interface{}) {
	toLog(ctx, columnsForCandidateLog, args)
}

func toLogCandidateVacancy(ctx *fasthttp.RequestCtx, args ...interface{}) {
	toLog(ctx, columnsForCandidateVacancyLog, args)
}

func toLogCompany(ctx *fasthttp.RequestCtx, args ...interface{}) {
	toLog(ctx, columnsForCompanyLog, args)
}

func toLogVacancy(ctx *fasthttp.RequestCtx, args ...interface{}) {
	toLog(ctx, columnsForVacancyLog, args)
}

func toLog(ctx *fasthttp.RequestCtx, columns []string, args []interface{}) {
	finalArgs := make([]interface{}, 0)
	finalArgs = append(finalArgs, auth.GetUserID(ctx), args, time.Now())

	_, err := db.LogsTable.Insert(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(finalArgs...))
	if err != nil {
		logs.ErrorLog(err, "toLog")
	}
}

func loLogUpdateValues(columns []string, args []interface{}) (ret map[string]interface{}) {
	if len(columns) > 0 {
		ret = make(map[string]interface{}, len(columns))
		for i, col := range columns {
			ret[col] = args[i]
		}
	}
	return
}
