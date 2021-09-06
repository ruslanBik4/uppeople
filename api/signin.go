// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"regexp"
	"time"

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

	opts := []dbEngine.BuildSqlOptions{
		dbEngine.WhereForSelect("email"),
		dbEngine.ArgsForSelect(a.Email),
	}
	users, _ := db.NewUsers(DB)
	err := users.SelectOneAndScan(ctx,
		u,
		opts...,
	)
	switch err {
	case nil:
		err := auth.CheckPass(u.Hash, a.Password)
		if err != nil {
			// req := &fasthttp.Request{}
			// ctx.Request.CopyTo(req)
			// if request(req, u) != nil {
			// 	logs.DebugLog(u)
			// 	return err.Error(), apis.ErrWrongParamsList
			// }
			// b, err := auth.NewHash(a.Password)
			// if err == nil {
			// 	_, err = users.Update(ctx,
			// 		dbEngine.ColumnsForSelect("password"),
			// 		dbEngine.WhereForSelect("id"),
			// 		dbEngine.ArgsForSelect(string(b), u.Id),
			// 	)
			// }
			// if err != nil {
			// 	logs.ErrorLog(err, "auth.NewHash")
			// }
			return createErrResult(err)
		}
	case pgx.ErrNoRows:
		return createErrResult(pgx.ErrNoRows)
	default:
		return createErrResult(err)
	}

	u.Token, err = auth.Bearer.NewToken(u)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNonAuthoritativeInfo)
		return nil, errors.Wrap(err, "Bearer.NewToken")
	}

	opts[1] = dbEngine.ArgsForSelect(time.Now(), getIP(ctx), a.Email)

	logs.DebugLog("login %s: %v, %s", a.Email, time.Now(), getIP(ctx))
	opts = append(opts,
		dbEngine.ColumnsForSelect("last_login", "last_ip"),
	)

	_, err = users.Update(ctx,
		opts...,
	)
	if err != nil {
		logs.ErrorLog(err, "users.Update")
	}

	return u, nil
}

var regIp = regexp.MustCompile(`for=s*(\d+\.?)+,`)

func getIP(ctx *fasthttp.RequestCtx) string {
	addr := ctx.Conn().RemoteAddr().String()
	if addr > "" {
		return addr
	}

	ipClient := ctx.Request.Header.Peek("X-Forwarded-For")
	addr = string(ipClient)
	if len(ipClient) == 0 {
		ipClient = ctx.Request.Header.Peek("Forwarded")
		ips := regIp.FindSubmatch(ipClient)

		if len(ips) == 0 {
			addr = string(ctx.Request.Header.Peek("X-ProxyUser-Ip"))
		} else {
			addr = string(ips[0])
		}
	}

	return addr
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
			u.RoleId = int32(user["role_id"].(float64))
			u.Phone, _ = user["tel"].(string)

			return nil
		}
	}

	return nil
}
