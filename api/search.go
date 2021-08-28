// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"
)

type GlobalSearch struct {
	Search string `json:"search"`
}

func (g *GlobalSearch) GetValue() interface{} {
	return g
}

func (g *GlobalSearch) NewValue() interface{} {
	return &GlobalSearch{}
}

func HandleGlobalSearch(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	dto, ok := ctx.UserValue(apis.JSONParams).(*GlobalSearch)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	var err error
	res := make(map[string][]map[string]interface{}, 0)
	where := " where name ~ $1 OR email ~$1 OR phone ~$1 OR skype ~$1 "

	res["candidates"], err = DB.Conn.SelectToMaps(ctx,
		"select * from candidates "+where,
		dto.Search,
	)
	if err != nil {
		return createErrResult(err)
	}

	res["companies"], err = DB.Conn.SelectToMaps(ctx,
		"select * from companies "+where,
		dto.Search,
	)
	if err != nil {
		return createErrResult(err)
	}

	res["contacts"], err = DB.Conn.SelectToMaps(ctx,
		"select * from contacts "+where,
		dto.Search,
	)
	if err != nil {
		return createErrResult(err)
	}

	res["vacancies"], err = DB.Conn.SelectToMaps(ctx,
		`select vacancies.id, vacancies.name as details, date_create as date, salary, l.name as location,
       platforms.name as platform, platforms.id as platId, seniorities.name as seniority,
       companies.id as companyId, companies.name as company, statuses.status, statuses.id as statusId
from vacancies left Join public.platforms on vacancies.platform_id=platforms.id
				left Join public.seniorities on vacancies.seniority_id = seniorities.id
    left Join companies on vacancies.company_id = companies.id
    left Join public.statuses on vacancies.status = statuses.id
    left Join public.location_for_vacancies l on vacancies.location_id = l.id
where vacancies.name ~ $1 OR l.name ~ $1 OR platforms.name ~ $1 OR seniorities.name ~ $1
   OR companies.name ~ $1
`,
		dto.Search,
	)
	if err != nil {
		return createErrResult(err)
	}

	return res, nil
}
