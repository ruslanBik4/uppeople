// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type UserResponse struct {
	Users       []*db.UsersFields `json:"users"`
	Partners    SelectedUnits     `json:"partners"`
	Freelancers SelectedUnits     `json:"freelancers"`
	Recruiters  SelectedUnits     `json:"recruiters"`
}

func HandleGetUser(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}
	users, _ := db.NewUsers(DB)
	err := users.SelectOneAndScan(ctx,
		users,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	return users.Record, nil
}

func HandleAllStaff(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	r := make([]*db.UsersFields, 0)
	users, _ := db.NewUsers(DB)
	err := users.SelectSelfScanEach(ctx,
		func(record *db.UsersFields) error {
			r = append(r, record)
			return nil
		})
	if err != nil {
		return createErrResult(err)
	}

	return UserResponse{
		Users: r,
	}, nil
}
