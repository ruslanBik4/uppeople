// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/valyala/fasthttp"
)

func HandleReturnLogsForCand(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return DB.Conn.SelectToMaps(ctx,
		`select logs.id as logId, logs.text as text, logs.d_c as date, 
		companies.id as compId, companies.name as compName, vacancies.id as vacId, 
		CONCAT_WS(' - ', platforms.nazva, seniorities.nazva) as vac
		from logs left Join companies on (logs.company_id = companies.id)
			left Join vacancies ON (logs.vacancy_id = vacancies.id)
			left Join platforms ON (vacancies.platform_id = platforms.id)
			left Join seniorities ON (vacancies.seniority_id = seniorities.id)
		where candidate_id =$1
		order by logs.d_c DESC`,
		ctx.UserValue("id"),
	)
}
