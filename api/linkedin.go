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
			Fnc:  HandleAuthlinkEdin,
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
			Fnc:  HandleGetPlatformslinkEdin,
			Desc: "get_platforms linkEdin",
		},
		"/api/get_selectors": {
			Fnc:  HandleGetSelectorslinkEdin,
			Desc: "get_selectors linkEdin",
		},
		"/api/get_seniorities": {
			Fnc:  HandleGetSenioritieslinkEdin,
			Desc: "get_seniorities linkEdin",
		},
		"/api/get_tags": {
			Fnc:  HandleGetTagslinkEdin,
			Desc: "get_tags linkEdin",
		},
		"/api/get_reasons": {
			Fnc:  HandleGetReasonslinkEdin,
			Desc: "get_reasons linkEdin",
		},
		"/api/get_recruiter_vacancies": {
			Fnc:  HandleGetRecruiterVacancieslinkEdin,
			Desc: "get_recruiter_vacancies linkEdin",
		},
		"/api/get_candidate_info": {
			Fnc:  HandleGetCandidate_infolinkEdin,
			Desc: "get_candidate_info linkEdin",
		},
	}
)

func HandleAuthlinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

func HandleGetPlatformslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_platforms",
	}

	return m, nil
}

func HandleGetCandidate_infolinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_candidate_info",
	}

	return m, nil
}

func HandleGetReasonslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_reasons",
	}

	return m, nil
}

func HandleGetRecruiterVacancieslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_recruiter_vacancies",
	}

	return m, nil
}

func HandleGetTagslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_tags",
	}

	return m, nil
}

func HandleGetSenioritieslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_seniorities",
	}

	return m, nil
}
func HandleGetSelectorslinkEdin(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
