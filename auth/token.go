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

type lastEdit struct {
	Candidate *db.CandidatesFields
	Vacancy   *db.VacanciesFields
}
type User struct {
	*db.UsersFields
	Companies map[int32]map[string]string `json:"companies"`
	Token     string                      `json:"token"`
	TokenOld  string                      `json:"token_old"`
	Host      string                      `json:"-"`
	LastEdit  lastEdit                    `json:"-"`
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

	logs.ErrorLog(dbEngine.ErrNotFoundColumn{}, "%s not user data but %T %[2]v", token)

	return nil
}

func PutEditCandidate(ctx *fasthttp.RequestCtx, data *db.CandidatesFields) bool {
	token, ok := ctx.UserValue(auth.UserValueToken).(*User)
	if ok {
		token.LastEdit.Candidate = data
	}

	return ok
}

func GetEditCandidate(ctx *fasthttp.RequestCtx) *db.CandidatesFields {
	token, ok := ctx.UserValue(auth.UserValueToken).(*User)
	if ok {
		return token.LastEdit.Candidate
	}

	return nil
}
