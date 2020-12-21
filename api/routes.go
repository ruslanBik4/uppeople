// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
)

var (
	ApiRoutes     = apis.NewMapRoutes()
	ApiPostRoutes = apis.ApiRoutes{
		"/api/auth/login": {
			Fnc:  HandleAuth,
			Desc: "show search results according range of characteristics",
			// DTO:    &DTOSearch{},
			// Resp:   search.RespGroups(),
		},
		"/api/main/addNewCandidate": {
			Fnc:      HandleAddCandidate,
			Desc:     "show search results according range of characteristics",
			NeedAuth: true,
		},
		// 'api/main/returnAllCandidates/':{
		// 	Fnc:  HandleAuth,
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

func HandleAuth(ctx *fasthttp.RequestCtx) (interface{}, error) {
	req := &fasthttp.Request{}
	ctx.Request.CopyTo(req)

	for _, host := range hosts {
		resp := &fasthttp.Response{}

		err := doRequest(req, resp, host)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode() == fasthttp.StatusOK {
			var v map[string]interface{}
			err := jsoniter.Unmarshal(resp.Body(), &v)
			u := &auth.User{
				Companies: make(map[int32]map[string]string),
				Host:      host,
				TokenOld:  v["access_token"].(string),
				// UsersFields: users.Record,
			}

			u.Token, err = auth.Bearer.NewToken(u)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusNonAuthoritativeInfo)
				return nil, errors.Wrap(err, "Bearer.NewToken")
			}

			v["access_token"] = u.Token

			return v, nil
		}
	}

	return nil, nil
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

func HandleAddCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	user, err := prepareRequest(ctx)
	if err != nil {
		return nil, err
	}

	var v map[string]interface{}
	err = jsoniter.Unmarshal(ctx.Request.Body(), &v)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	req := &fasthttp.Request{}
	ctx.Request.CopyTo(req)

	req.URI().SetPath("api/main/allCandidates/1")
	nameCan := v["name"].(string)
	req.SetBodyString(`{"search":"` + nameCan + `","dateFrom":"","dateTo":"","selectCompanies":[],"selectPlatforms":[],"selectStatuses":[],"selectRecruiter":"","selectTag":[],"selectReason":[],"mySent":false,"dateFromAllCandidates":"","dateToAllCandidates":"","dateFromSentCandidates":"","dateToSentCandidates":"","dateFromFreelancersCandidates":"","dateToFreelancersCandidates":"","dateFollowUpFrom":"","dateFollowUpTo":"","selectSeniority":[]}`)
	resp := &fasthttp.Response{}
	err = doRequest(req, resp, user.Host)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(resp.Header.ContentType(), []byte("json")) {
		err = jsoniter.Unmarshal(resp.Body(), &v)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal")
		}

		c, ok := v["Count"].(float64)
		if !ok {
			logs.DebugLog("%T", v["Count"])
		} else if c > 0 {
			return nameCan + " already exists", apis.ErrWrongParamsList
		} else {
			logs.DebugLog("%f", v["Count"])
			return nil, doRequest(&ctx.Request, &ctx.Response, user.Host)
		}
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
