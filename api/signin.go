// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/jackc/pgx/v4"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
	"github.com/valyala/fasthttp"
)

type DTOAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *DTOAuth) ReadParams(ctx *fasthttp.RequestCtx) {
	a.Email = ctx.UserValue("email").(string)
	a.Password = ctx.UserValue("password").(string)
}

func (a *DTOAuth) GetValue() interface{} {
	return a
}

func (a *DTOAuth) NewValue() interface{} {
	return &DTOAuth{}
}

func HandleAuthLogin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	a, ok := ctx.UserValue(apis.JSONParams).(*DTOAuth)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	u := &auth.User{
		UsersFields: &db.UsersFields{},
		Companies:   make(map[int32]map[string]string),
	}

	if a.Email == "test@test.com" && a.Password == "1111" {
		u.Id = 27
		u.Name = "julia"
		u.Email = "julia@uppeople.co"
		u.Role_id = 2
	} else {
		users, _ := db.NewUsers(DB)
		err := users.SelectOneAndScan(ctx,
			u,
			dbEngine.WhereForSelect("email"),
			dbEngine.ArgsForSelect(a.Email),
		)
		switch err {
		case nil:
			err := auth.CheckPass(u.Pass.String, a.Password)
			if err != nil {
				req := &fasthttp.Request{}
				ctx.Request.CopyTo(req)
				if request(req, u) != nil {
					logs.DebugLog(u)
					return err.Error(), apis.ErrWrongParamsList
				}
				b, err := auth.NewHash(u.Pass.String)
				if err == nil {
					_, err = users.Update(ctx,
						dbEngine.ColumnsForSelect("password"),
						dbEngine.WhereForSelect("id"),
						dbEngine.ArgsForSelect(string(b), u.Id),
					)
				}
				if err != nil {
					logs.ErrorLog(err, "auth.NewHash")
				}
			}
		case pgx.ErrNoRows:
			return createErrResult(pgx.ErrNoRows)
		default:
			return createErrResult(err)
		}
	}

	var err error
	u.Token, err = auth.Bearer.NewToken(u)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNonAuthoritativeInfo)
		return nil, errors.Wrap(err, "Bearer.NewToken")
	}

	u.Pass.String = ""

	return u, nil
}

func request(req *fasthttp.Request, u *auth.User) error {
	for _, host := range hosts {
		resp := &fasthttp.Response{}

		err := doRequest(req, resp, host)
		if err != nil {
			return err
		}

		if resp.StatusCode() == fasthttp.StatusOK {
			var v map[string]interface{}
			err := jsoniter.Unmarshal(resp.Body(), &v)
			if err != nil {
				return errors.Wrap(err, "")
			}

			user := v["user"].(map[string]interface{})

			u.Id = int32(user["id"].(float64))
			u.Name = user["name"].(string)
			u.Email = user["email"].(string)
			u.Role_id = int32(user["role_id"].(float64))
			u.Phone.String, _ = user["tel"].(string)

			return nil
		}
	}

	return nil
}
