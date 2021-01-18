// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"go/types"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
)

var (
	ParamID = apis.InParam{
		Name:     "id",
		Desc:     "id of candidate",
		DefValue: apis.ApisValues(apis.ChildRoutePath),
		Req:      true,
		Type:     apis.NewTypeInParam(types.Int32),
	}
)

var (
	Routes     = apis.NewMapRoutes()
	PostRoutes = apis.ApiRoutes{
		"/api/auth/login": {
			Fnc:  HandleAuthLogin,
			Desc: "show search results according range of characteristics",
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
			// Resp:   search.RespGroups(),
		},
		"/api/main/addNewVacancy": {
			Fnc:      HandleAddVacancy,
			Desc:     "add new vacancy",
			DTO:      &VacancyDTO{},
			NeedAuth: true,
		},
		"/api/main/editVacancy/": {
			Fnc:      HandleEditVacancy,
			Desc:     "edit vacancy",
			DTO:      &VacancyDTO{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/addNewCandidate": {
			Fnc:      HandleAddCandidate,
			Desc:     "show search results according range of characteristics",
			DTO:      &CandidateDTO{},
			NeedAuth: true,
		},
		"/api/main/editCandidate/": {
			Fnc:      HandleEditCandidate,
			Desc:     "show search results according range of characteristics",
			DTO:      &CandidateDTO{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/saveFollowUp": {
			Fnc:      HandleFollowUpCandidate,
			Desc:     "show search results according range of characteristics",
			DTO:      &FollowUpDTO{},
			NeedAuth: true,
		},
		// 'api/main/returnAllCandidates/':{
		// 	Fnc:  HandleAuthLogin,
		// 	Desc: "show search results according range of characteristics",
		// },
		"/api/main/viewAllVacancyInCompany/": {
			Fnc:  HandleViewAllVacancyInCompany,
			Desc: "get list of vacancies",
			DTO:  &vacDTO{},
		},
		"/api/": {
			Fnc:  HandleApiRedirect,
			Desc: "show search results according range of characteristics",
			// DTO:    &DTOSearch{},
			Method:   apis.POST,
			NeedAuth: true,
			// Resp:   search.RespGroups(),
		},
		"api/main/getTags": {
			Fnc:      HandleGetTags,
			Desc:     "return select of tags",
			Method:   apis.POST,
			NeedAuth: true,
		},
		"api/main/getStatuses": {
			Fnc:      HandleGetStatuses,
			Desc:     "return select of tags",
			Method:   apis.POST,
			NeedAuth: true,
		},
		"/api/main/allCandidates/": {
			Fnc:      HandleAllCandidate,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &SearchCandidates{},
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/getCandidatesAmountByTags": {
			Fnc:      HandleGetCandidatesAmountByTags,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &DTOAmounts{},
		},
	}
	GetRoutes = apis.ApiRoutes{
		"/api/admin/all-staff": {
			Fnc:      HandleAllStaff,
			Desc:     "show all candidates",
			NeedAuth: true,
		},
		"/api/main/returnOptionsForSelects": {
			Fnc:      HandleReturnOptionsForSelects,
			Desc:     "show all candidates",
			NeedAuth: true,
		},
		"/api/main/allCandidates/": {
			Fnc:      HandleAllCandidate,
			Desc:     "show all candidates",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/viewOneCandidate/": {
			Fnc:      HandleViewCandidate,
			Desc:     "show one candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/viewVacancy/": {
			Fnc:      HandleViewVacancy,
			Desc:     "show one candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/admin/returnLogsForCand/": {
			Fnc:      HandleReturnLogsForCand,
			Desc:     "show logs of candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		// "/api/main/returnAllCandidates/": {
		// 	Fnc:      HandleAllCandidate,
		// 	Desc:     "show search results according range of characteristics",
		// 	NeedAuth: true,
		// },
		// "/api/main/viewCandidatesFreelancerOnVacancies/": {
		// 	Fnc:      HandleAllCandidate,
		// 	Desc:     "show search results according range of characteristics",
		// 	NeedAuth: true,
		// },
		"/api/": {
			Fnc:      HandleApiRedirect,
			Desc:     "show search results according range of characteristics",
			NeedAuth: true,
			// DTO:    &DTOSearch{},
			// Method: apis.POST,
			// Resp:   search.RespGroups(),
		},
	}
)
var hosts = []string{
	"back.uppeople.co",
	"back2.uppeople.co",
}

func doRequest(req *fasthttp.Request, resp *fasthttp.Response, host string) error {
	c := fasthttp.Client{}
	uri := req.URI()
	uri.SetScheme("http")
	uri.SetHost(host)

	err := c.Do(req, resp)
	if err != nil {
		return errors.Wrap(err, "do")
	}

	return nil
}

func HandleApiRedirect(ctx *fasthttp.RequestCtx) (interface{}, error) {
	user, err := prepareRequest(ctx)
	if err != nil {
		return nil, err
	}

	req := &fasthttp.Request{}
	ctx.Request.CopyTo(req)

	err = doRequest(req, &ctx.Response, user.Host)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func prepareRequest(ctx *fasthttp.RequestCtx) (*auth.User, error) {
	user := auth.GetUserData(ctx)
	if user == nil {
		return nil, apis.ErrWrongParamsList
	}

	ctx.Request.Header.Set("Authorization", "Bearer "+user.TokenOld)
	return user, nil
}

func init() {
	for _, route := range PostRoutes {
		route.Method = apis.POST
	}
	Routes.AddRoutes(PostRoutes)
	Routes.AddRoutes(GetRoutes)
}
