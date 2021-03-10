// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"go/types"
	"strings"

	"github.com/jackc/pgx/v4"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/uppeople/auth"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

var (
	LERoutes = apis.ApiRoutes{
		"/api/authorize": {
			Fnc:  HandleAuthLinkedin,
			Desc: "auth for linkEdin",
			Params: []apis.InParam{
				{
					Name: "email",
					Type: apis.NewTypeInParam(types.String),
				},
				{
					Name: "password",
					Type: apis.NewTypeInParam(types.String),
				},
			},
			DTO: &DTOAuth{},
		},
		"/api/get_platforms": {
			Fnc:  HandleGetPlatformsLinkedin,
			Desc: "get_platforms linkEdin",
		},
		"/api/get_selectors": {
			Fnc:  HandleGetSelectorslinkEdin,
			Desc: "get_selectors linkEdin",
		},
		"/api/get_seniorities": {
			Fnc:  HandleGetSenioritiesLinkedin,
			Desc: "get_seniorities linkEdin",
		},
		"/api/get_tags": {
			Fnc:  HandleGetTagsLinkedin,
			Desc: "get_tags linkEdin",
		},
		"/api/get_reasons": {
			Fnc:  HandleGetReasonsLinkedin,
			Desc: "get_reasons linkEdin",
		},
		"/api/get_recruiter_vacancies": {
			Fnc:      HandleGetRecruiterVacancieslinkEdin,
			Desc:     "get_recruiter_vacancies linkEdin",
			DTO:      &DTOVacancy{},
			NeedAuth: true,
			WithCors: true,
		},
		"/api/get_candidate_info": {
			Fnc:      HandleGetCandidate_infolinkEdin,
			Desc:     "get_candidate_info linkEdin",
			WithCors: true,
			Params: []apis.InParam{
				{
					Name: "url",
					Type: apis.NewTypeInParam(types.String),
				},
			},
		},
	}
)

func HandleAuthLinkedin(ctx *fasthttp.RequestCtx) (interface{}, error) {

	str, err := jsoniter.Marshal(ctx.UserValue(apis.JSONParams))
	if err != nil {
		return createErrResult(errors.New("wrong DTO"))
	}

	ctx.Request.ResetBody()
	ctx.Request.Header.SetContentType("application/json; charset=utf-8")
	ctx.Request.SetBody(str)

	return HandleAuthLogin(ctx)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "")
	// }
	//
	//
	// m := map[string]interface{}{
	// 	"status":  "success",
	// 	"user_id": v.(map[string]interface{})["user"].(map[string]interface{})["id"],
	// }
	//
	// return m, nil
}

func HandleGetPlatformsLinkedin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	m := map[string]interface{}{
		"status": getPlatforms(ctx, DB),
	}

	return m, nil
}

func HandleGetCandidate_infolinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	url, ok := ctx.UserValue("url").(string)
	if !ok {
		return map[string]interface{}{
			"url": "wrong type",
		}, apis.ErrWrongParamsList
	}
	table, _ := db.NewCandidates(DB)
	err := table.SelectOneAndScan(ctx, table, dbEngine.WhereForSelect("linkedin"), dbEngine.ArgsForSelect(url))
	if err != nil && err != pgx.ErrNoRows {
		return createErrResult(err)
	}

	m := map[string]interface{}{
		"status": "ok",
		"data":   []string{},
	}

	if err == nil {
		id := table.Record.Id
		m["crm_url"] = fmt.Sprintf("%s://%s/#/candidates/%d", ctx.URI().Scheme(), ctx.Host(), id)
		m["data"] = table.Record
	}

	return m, nil
}

func HandleGetReasonsLinkedin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	m := map[string]interface{}{
		"status": getRejectReason(ctx, DB),
	}

	return m, nil
}

func HandleGetRecruiterVacancieslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	data, err := DB.Conn.SelectToMaps(ctx,
		`select vacancies.id, platform_id, user_ids,
CONCAT(companies.name, ' (', platforms.nazva , ')') as name,
LOWER(CONCAT(companies.name, ' (', platforms.nazva , ')')) as label
from vacancies left join companies on vacancies.company_id = companies.id
left join platforms on vacancies.platform_id = platforms.id
where $1=ANY(user_ids)`,
		auth.GetUserData(ctx).Id,
	)
	if err != nil {
		return createErrResult(err)
	}

	m := map[string]interface{}{
		"status":    "ok",
		"vacancies": data,
	}

	return m, nil
}

func HandleGetTagsLinkedin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	m := map[string]interface{}{
		"status": getTags(ctx, DB),
	}

	return m, nil
}

func HandleGetSenioritiesLinkedin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	m := map[string]interface{}{
		"status": getSeniorities(ctx, DB),
	}

	return m, nil
}
func HandleGetSelectorslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewChrome_extention_selectors(DB)
	records := make([]*db.Chrome_extention_selectorsFields, 0)
	err := table.SelectSelfScanEach(ctx,
		func(record *db.Chrome_extention_selectorsFields) error {
			records = append(records, record)
			return nil
		})
	if err != nil {
		return createErrResult(err)
	}
	m := map[string]interface{}{
		"status": "ok",
		"data":   records,
	}

	return m, nil
}

func init() {
	for path, route := range LERoutes {
		route.WithCors = true
		if !strings.HasSuffix(path, "show_candidate_linkedin") {
			route.Method = apis.POST
		}
	}
	Routes.AddRoutes(LERoutes)
}
