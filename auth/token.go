// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/auth"
	"github.com/ruslanBik4/logs"
	"github.com/ruslanBik4/uppeople/db"
	"github.com/valyala/fasthttp"
)

type lastEdit struct {
	Candidate *db.CandidatesFields
	Vacancy   *db.VacanciesFields
}

type User struct {
	*db.UsersFields
	Companies map[int32]map[string]string `json:"companies"`
	Pass      sql.NullString              `json:"password"`
	Token     string                      `json:"access_token"`
	LastEdit  lastEdit                    `json:"-"`
}

func (u *User) GetFields(columns []dbEngine.Column) []interface{} {
	values := make([]interface{}, len(columns))
	for i, col := range columns {
		switch name := col.Name(); name {
		case "password":
			values[i] = &u.Pass
		default:
			values[i] = u.RefColValue(col.Name())
		}
	}

	return values
}

func (u *User) IsAdmin() bool {
	return u.RoleId == 1
}

func (u *User) GetUserID() int {
	return int(u.Id)
}

func (u *User) GetSchema() string {
	return u.Schema
}

func GetUserID(ctx *fasthttp.RequestCtx) int {
	token, ok := ctx.UserValue(auth.UserValueToken).(*User)
	if ok {
		return token.GetUserID()
	}

	logs.ErrorLog(dbEngine.ErrNotFoundColumn{}, "%s not user data but %T %[2]v", token)

	return -1
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

func PutEditVacancy(ctx *fasthttp.RequestCtx, data *db.VacanciesFields) bool {
	token, ok := ctx.UserValue(auth.UserValueToken).(*User)
	if ok {
		token.LastEdit.Vacancy = data
	}

	return ok
}

func GetEditVacancy(ctx *fasthttp.RequestCtx) *db.VacanciesFields {
	token, ok := ctx.UserValue(auth.UserValueToken).(*User)
	if ok {
		return token.LastEdit.Vacancy
	}

	return nil
}
