// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"go/types"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"
)

var (
	LERoutes = apis.ApiRoutes{
		"/api/authorize": {
			Fnc:  HandleAuthLinkEdit,
			Desc: "auth for linkEdit",
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
			Fnc:  HandleGetPlatformsLinkEdit,
			Desc: "get_platforms LInkEdit",
		},
		"/api/get_selectors": {
			Fnc:  HandleGetSelectorsLinkEdit,
			Desc: "get_selectors LInkEdit",
		},
		"/api/get_seniorities": {
			Fnc:  HandleGetSenioritiesLinkEdit,
			Desc: "get_seniorities LInkEdit",
		},
		"/api/get_tags": {
			Fnc:  HandleGetTagsLinkEdit,
			Desc: "get_tags LInkEdit",
		},
		"/api/get_reasons": {
			Fnc:  HandleGetReasonsLinkEdit,
			Desc: "get_reasons LInkEdit",
		},
		"/api/get_recruiter_vacancies": {
			Fnc:  HandleGetRecruiterVacanciesLinkEdit,
			Desc: "get_recruiter_vacancies LInkEdit",
		},
		"/api/get_candidate_info": {
			Fnc:  HandleGetCandidate_infoLinkEdit,
			Desc: "get_candidate_info LInkEdit",
		},
	}
)

func HandleAuthLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	v, err := HandleAuthLogin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	m := map[string]interface{}{
		"status":  "success",
		"user_id": v.(map[string]interface{})["user"].(map[string]interface{})["id"],
	}

	return m, nil
}

func HandleGetPlatformsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_platforms",
	}

	return m, nil
}

func HandleGetCandidate_infoLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_candidate_info",
	}

	return m, nil
}

func HandleGetReasonsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_reasons",
	}

	return m, nil
}

func HandleGetRecruiterVacanciesLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_recruiter_vacancies",
	}

	return m, nil
}

func HandleGetTagsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_tags",
	}

	return m, nil
}

func HandleGetSenioritiesLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_seniorities",
	}

	return m, nil
}
func HandleGetSelectorsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_selectors",
	}

	return m, nil
}

func init() {
	for _, route := range LERoutes {
		route.WithCors = true
	}
	Routes.AddRoutes(LERoutes)
}
