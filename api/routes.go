// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

var (
	ApiRoutes   = apis.NewMapRoutes()
	AdminRoutes = apis.ApiRoutes{
		"/api/": {
			Fnc:  HandleApiRedirect,
			Desc: "show search results according range of characteristics",
			// DTO:    &DTOSearch{},
			Method: apis.POST,
			// Resp:   search.RespGroups(),
		},
	}
	SearchRoutes = apis.ApiRoutes{
		"/api/": {
			Fnc:  HandleApiRedirect,
			Desc: "show search results according range of characteristics",
			// DTO:    &DTOSearch{},
			// Method: apis.POST,
			// Resp:   search.RespGroups(),
		},
	}
)

func HandleApiRedirect(ctx *fasthttp.RequestCtx) (interface{}, error) {
	uri := ctx.Request.URI()
	uri.SetHost("back.uppeople.co")
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
