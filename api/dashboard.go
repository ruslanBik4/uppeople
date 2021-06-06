// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/valyala/fasthttp"
)

func HandleDashBoard(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	// r, ok := DB.Routines["dashboard"]
	// if !ok {
	// 	ctx.SetStatusCode(fasthttp.StatusNoContent)
	// 	return nil, nil
	// }

	return DB.Conn.SelectToMap(ctx,
		`select * from dashboard()`)
}
