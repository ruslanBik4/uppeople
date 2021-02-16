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

	return DB.Conn.SelectToMaps(ctx,
		`select logs.id as logId, CONCAT('Пользователь ', users.name, 
		CASE WHEN kod_deystviya=102 THEN ' проработал ' 
			 WHEN kod_deystviya=101  THEN ' добавил нового '
			 WHEN kod_deystviya=100  THEN ' обновил у '
			 WHEN kod_deystviya=103  THEN ' удалил '
			ELSE '' END,
		CASE WHEN candidate_id > 0 THEN CONCAT(' кандидата ', can.name)
			 WHEN vacancy_id > 0 THEN CONCAT(' вакансию компании ', companies.name)
			ELSE '' END,
			' ', logs.text) as text, 
		logs.d_c as date, 
		companies.id as compId, companies.name as compName, vacancies.id as vacId, 
		CONCAT_WS(' - ', platforms.nazva, seniorities.nazva) as vac
		from logs left Join companies on (logs.company_id = companies.id)
			left join vacancies ON (logs.vacancy_id = vacancies.id)
			join users ON (logs.user_id = users.id)
			join candidates can ON (logs.candidate_id = can.id)
			left Join platforms ON (vacancies.platform_id = platforms.id)
			left Join seniorities ON (vacancies.seniority_id = seniorities.id)
		where candidate_id =$1
		order by logs.d_c DESC`,
		ctx.UserValue("id"),
	)
}

func toLogCandidate(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId int32, text string, code int32) {
	user := auth.GetUserData(ctx)
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "text", "date_create", "d_c",
			"kod_deystviya"),
		dbEngine.ArgsForSelect(user.Id, candidateId,
			text,
			time.Now(),
			time.Now(),
			code))
}

func toLogCandidateVacancy(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, candidateId, companyId, vacancyId int32, text string, code int32) {
	user := auth.GetUserData(ctx)
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "company_id", "vacancy_id", "text", "date_create", "d_c",
			"kod_deystviya"),
		dbEngine.ArgsForSelect(user.Id, candidateId, companyId, vacancyId,
			text,
			time.Now(),
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
