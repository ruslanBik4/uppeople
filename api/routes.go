// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
)

var (
	ApiRoutes   = apis.NewMapRoutes()
	AdminRoutes = apis.ApiRoutes{
		"/api/auth/login": {
			Fnc:  HandleAuth,
			Desc: "show search results according range of characteristics",
			// DTO:    &DTOSearch{},
			Method: apis.POST,
			// Resp:   search.RespGroups(),
		},
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
	c := fasthttp.Client{}

	for _, host := range hosts {
		uri := ctx.Request.URI()
		uri.SetScheme("http")
		uri.SetHost(host)
		reg := fasthttp.Request{}
		ctx.Request.CopyTo(&reg)

		resp := &fasthttp.Response{}
		err := c.Do(&reg, resp)
		if err != nil {
			return nil, errors.Wrap(err, "do")
		}

		if resp.StatusCode() != fasthttp.StatusUnauthorized {
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

			if v["access_token"].(string) != u.Token {
				logs.DebugLog("token: '%s' \n'%s'", v["access_token"].(string), u.Token)
			}
			v["access_token"] = u.Token
			return v, nil
		}
	}

	return nil, nil
}

func HandleApiRedirect(ctx *fasthttp.RequestCtx) (interface{}, error) {
	user := auth.GetUserData(ctx)
	if user == nil {
		return nil, apis.ErrWrongParamsList
	}
	uri := ctx.Request.URI()
	uri.SetHost(user.Host)
	ctx.Request.Header.Set("access_token", user.TokenOld)
	// uri.SetScheme("http")
	// reg := fasthttp.Request{}
	// ctx.Request.CopyTo(&reg)
	//
	// ctx.Request = reg
	ctx.RedirectBytes(uri.FullURI(), fasthttp.StatusPermanentRedirect)
	logs.DebugLog("redirect %s %s", string(uri.FullURI()), ctx.Method())

	return nil, nil
}

func init() {
	ApiRoutes.AddRoutes(AdminRoutes)
	ApiRoutes.AddRoutes(SearchRoutes)
}
