// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/auth"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type User struct {
	*db.UsersFields
	Companies map[int32]map[string]string `json:"companies"`
	Token     string                      `json:"token"`
	TokenOld  string                      `json:"token"`
	Host      string                      `json:"-"`
}

func (u *User) IsAdmin() bool {
	return u.Id_roles == 1
}

func (u *User) GetUserID() int {
	return int(u.Id)
}

func GetUserData(ctx *fasthttp.RequestCtx) *User {
	token, ok := ctx.UserValue(auth.UserValueToken).(*User)
	if ok {
		return token
	}

	logs.ErrorLog(dbEngine.ErrNotFoundColumn{}, "not user data but %T", token)

	return nil
}
