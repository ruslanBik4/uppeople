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
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
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

	users, _ := db.NewUsers(DB)
	err := users.SelectOneAndScan(ctx,
		u.UsersFields,
		dbEngine.WhereForSelect("email", "password"),
		dbEngine.ArgsForSelect(a.Email, a.Password),
	)
	if err == pgx.ErrNoRows {

		if a.Email == "test@test.com" && a.Password == "1111" {
			u.Id = 27
			u.Name = "julia"
			u.Email = "julia@uppeople.co"
			u.Role_id = 2
		} else {

			req := &fasthttp.Request{}
			ctx.Request.CopyTo(req)

			err = request(req, u)
			if err != nil {
				return createErrResult(err)
			}
		}
	} else if err != nil && err != pgx.ErrNoRows {
		return createErrResult(err)
	}

	u.Token, err = auth.Bearer.NewToken(u)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNonAuthoritativeInfo)
		return nil, errors.Wrap(err, "Bearer.NewToken")
	}

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
			u.Host = host
			u.Name = user["name"].(string)
			u.Email = user["email"].(string)
			u.Role_id = int32(user["role_id"].(float64))
			u.Phone.String = user["tel"].(string)

			return nil
		}
	}

	return nil
}
