// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/httpgo/apis"
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
	SearchRoutes = apis.ApiRoutes{}
)

func HandleApiRedirect(ctx *fasthttp.RequestCtx) (interface{}, error) {
	ctx.Redirect("http://back2.uppeople.co"+string(ctx.URI().Path()), fasthttp.StatusMovedPermanently)
	return nil, nil
}

func init() {
	ApiRoutes.AddRoutes(AdminRoutes)
	ApiRoutes.AddRoutes(SearchRoutes)
}
