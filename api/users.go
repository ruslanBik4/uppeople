// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type UserResponse struct {
	Users       []*db.UsersFields `json:"users"`
	Partners    SelectedUnits     `json:"partners"`
	Freelancers SelectedUnits     `json:"freelancers"`
	Recruiters  SelectedUnits     `json:"recruiters"`
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
		return nil, errors.Wrap(err, "users.SelectSelfScanEach")
	}
	return UserResponse{
		Users: r,
	}, nil
}
