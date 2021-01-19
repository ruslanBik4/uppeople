// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/valyala/fasthttp"
)

type UserResponse struct {
	Users       SelectedUnits `json:"users"`
	Partners    SelectedUnits `json:"partners"`
	Freelancers SelectedUnits `json:"freelancers"`
	Recruiters  SelectedUnits `json:"recruiters"`
}

func HandleAllStaff(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return UserResponse{
		Users: getRecruters(ctx, DB),
	}, nil
}
