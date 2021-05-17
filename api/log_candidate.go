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

	return DB.Conn.SelectToMaps(ctx,
		fmt.Sprintf(`select logs.id as logId, CONCAT('Пользователь ', users.name, 
		CASE WHEN kod_deystviya=%d THEN ' проработал ' 
			 WHEN kod_deystviya=%d  THEN ' добавил нового '
			 WHEN kod_deystviya=%d  THEN ' обновил у '
			 WHEN kod_deystviya=%d  THEN ' удалил '
			ELSE '' END,
		CASE WHEN candidate_id > 0 THEN CONCAT(' кандидата ', can.name)
			 WHEN vacancy_id > 0 THEN CONCAT(' вакансию компании ', companies.name)
			ELSE '' END,
			' ', logs.text) as text, 
		logs.create_at as date, 
		companies.id as compId, companies.name as compName, vacancies.id as vacId, 
		CONCAT_WS(' - ', platforms.name, seniorities.name) as vac
		from logs left Join companies on (logs.company_id = companies.id)
			left join vacancies ON (logs.vacancy_id = vacancies.id)
			join users ON (logs.user_id = users.id)
			join candidates can ON (logs.candidate_id = can.id)
			left Join platforms ON (vacancies.platform_id = platforms.id)
			left Join seniorities ON (vacancies.seniority_id = seniorities.id)
		where candidate_id =$1
		order by logs.create_at DESC`,
			CODE_LOG_PEFORM, CODE_LOG_INSERT, CODE_LOG_UPDATE, CODE_LOG_DELETE),
		ctx.UserValue("id"),
	)
}

func toLogCandidate(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "text", "date_create",
			"kod_deystviya"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), candidateId,
			text,
			time.Now(),
			code))
}

func toLogCandidateVacancy(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId, companyId, vacancyId int32, text string, code int32) {
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "company_id", "vacancy_id", "text", "date_create",
			"kod_deystviya"),
		dbEngine.ArgsForSelect(auth.GetUserID(ctx), candidateId, companyId, vacancyId,
			text,
			time.Now(),
			code))
}

func toLog(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, columns, args dbEngine.BuildSqlOptions) {
	log, _ := db.NewLogs(DB)
	_, err := log.Insert(ctx, columns, args)
	if err != nil {
		logs.ErrorLog(err, "toLogCandidate")
	}
}
