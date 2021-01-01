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
	ApiRoutes     = apis.NewMapRoutes()
	ApiPostRoutes = apis.ApiRoutes{
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
			// DTO:    &DTOSearch{},
			// Resp:   search.RespGroups(),
		},
		"/api/main/addNewCandidate": {
			Fnc:      HandleAddCandidate,
			Desc:     "show search results according range of characteristics",
			DTO:      &CandidateDTO{},
			NeedAuth: true,
		},
		// 'api/main/returnAllCandidates/':{
		// 	Fnc:  HandleAuthLogin,
		// 	Desc: "show search results according range of characteristics",
		// },
		"/api/": {
			Fnc:  HandleApiRedirect,
			Desc: "show search results according range of characteristics",
			// DTO:    &DTOSearch{},
			Method:   apis.POST,
			NeedAuth: true,
			// Resp:   search.RespGroups(),
		},
	}
	SearchRoutes = apis.ApiRoutes{
		"/api/main/allCandidates/": {
			Fnc:      HandleAllCandidate,
			Desc:     "show all candidates",
			NeedAuth: true,
		},
		"/api/main/viewOneCandidate/": {
			Fnc:      HandleViewCandidate,
			Desc:     "show one candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				{
					Name:     "id",
					Desc:     "id of candidate",
					DefValue: apis.ApisValues(apis.ChildRoutePath),
					Req:      true,
					Type:     apis.NewTypeInParam(types.Int32)},
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
	for _, route := range ApiPostRoutes {
		route.Method = apis.POST
	}
	ApiRoutes.AddRoutes(ApiPostRoutes)
	ApiRoutes.AddRoutes(SearchRoutes)
}
