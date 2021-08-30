// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

var (
	columnsForCandidateLog        = []string{"candidate_id", "action_code", "changed"}
	columnsForCandidateVacancyLog = []string{"candidate_id", "company_id", "vacancy_id", "action_code", "changed"}
	columnsForVacancyLog          = []string{"company_id", "vacancy_id", "action_code", "changed"}
	columnsForCompanyLog          = []string{"company_id", "action_code", "changed"}
	columnsForCandidateStatusLog  = []string{"candidate_id", "vacancy_id", "action_code", "changed"}
)

func toLogCandidateUpdate(ctx *fasthttp.RequestCtx, candidateId int32, text map[string]interface{}) {
	toLogCandidate(ctx, candidateId, db.GetLogUpdateId(), text)
}

func toLogCandidatePerform(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, db.GetLogPerformId(), text)
}

func toLogCandidateDelete(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, db.GetLogDeleteId(), text)
}

func toLogCandidateRecontact(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, db.GetLogReContactId(), text)
}

func toLogCandidateAddComment(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, db.GetLogAddCommentId(), text)
}

func toLogCandidateDelComment(ctx *fasthttp.RequestCtx, candidateId int32, text string) {
	toLogCandidate(ctx, candidateId, db.GetLogDelCommentId(), text)
}

func toLogCandidateSendCV(ctx *fasthttp.RequestCtx, candidateId, companyId, vacancyId int32, text string) {
	toLogCandidateVacancy(ctx, candidateId, companyId, vacancyId, db.GetLogSendCVId(), text)
}

func toLogCandidateAppointInterview(ctx *fasthttp.RequestCtx, candidateId, companyId, vacancyId int32, text string) {
	toLogCandidateVacancy(ctx, candidateId, companyId, vacancyId, db.GetLogAppointInterviewId(), text)
}

func toLogCompanyInsert(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, db.GetLogInsertId(), text)
}

func toLogCompanyUpdate(ctx *fasthttp.RequestCtx, companyId int32, text map[string]interface{}) {
	toLogCompany(ctx, companyId, db.GetLogUpdateId(), text)
}

func toLogCompanyAddComment(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, db.GetLogAddCommentId(), text)
}

func toLogCompanyDelComment(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, db.GetLogDelCommentId(), text)
}

func toLogCompanyDelete(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, db.GetLogDeleteId(), text)
}

func toLogCompanyPerform(ctx *fasthttp.RequestCtx, companyId int32, text string) {
	toLogCompany(ctx, companyId, db.GetLogPerformId(), text)
}

func toLogVacancyInsert(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, companyId, vacancyId, db.GetLogInsertId(), text)
}

func toLogVacancyUpdate(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text map[string]interface{}) {
	toLogVacancy(ctx, companyId, vacancyId, db.GetLogUpdateId(), text)
}

func toLogVacancyPerform(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, companyId, vacancyId, db.GetLogPerformId(), text)
}

func toLogVacancyDelete(ctx *fasthttp.RequestCtx, companyId, vacancyId int32, text string) {
	toLogVacancy(ctx, companyId, vacancyId, db.GetLogDeleteId(), text)
}

func toLogCandidateUpdateStatus(ctx *fasthttp.RequestCtx, candidateId, vacancyId int32, text map[string]interface{}) {
	args := []interface{}{candidateId, vacancyId, db.GetLogUpdateId(), text}
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
	if val, ok := args[len(args)-1].(string); ok {
		args = append(args[:len(args)-1], map[string]string{"text": val})
	}

	// todo: will decide about ctx params which will be added to content
	const tokenName = auth.ValueTokenName
	user, ok := ctx.Value(tokenName).(*auth.User)
	if !ok {
		logs.ErrorLog(dbEngine.ErrWrongType{
			Name:     tokenName,
			TypeName: fmt.Sprintf("%T", ctx.Value(tokenName)),
		}, "can't write to Log")
		return
	}

	args = append(args, auth.GetUserID(ctx))
	columns = append(columns, "user_id")

	ctxTodo := context.WithValue(context.TODO(), tokenName, user)
	go db.InsertToLog(ctxTodo,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(args...))
}

func toLogUpdateValues(columns []string, args []interface{}) (ret map[string]interface{}) {
	if len(columns) > 0 {
		ret = make(map[string]interface{}, len(columns))
		for i, col := range columns {
			ret[col] = args[i]
		}
	}
	return
}
