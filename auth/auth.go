// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/httpgo/auth"
	"github.com/ruslanBik4/httpgo/services"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

// auth block
var (
	tokens  = auth.NewMapTokens(time.Hour * 24)
	Bearer  = auth.NewAuthBearer(tokens)
	Basic   = auth.NewAuthBasic(tokens, NewBasic)
	AuthMix = NewMixAuth(Bearer, Basic)
)

type MixAuth struct {
	auths []apis.FncAuth
}

func NewMixAuth(auth ...apis.FncAuth) *MixAuth {
	return &MixAuth{auth}
}

func (a *MixAuth) Auth(ctx *fasthttp.RequestCtx) bool {
	for _, a := range a.auths {
		if a.Auth(ctx) {
			return true
		}
	}

	return false
}

func (a *MixAuth) AdminAuth(ctx *fasthttp.RequestCtx) bool {
	for _, a := range a.auths {
		if a.AdminAuth(ctx) {
			return true
		}
	}

	return false
}
func (a *MixAuth) String() string {
	s := "MixAuth:"
	for _, a := range a.auths {
		s += a.String() + " "
	}

	return s
}

func NewBasic(ctx *fasthttp.RequestCtx, user, pass []byte) auth.TokenData {
	hash := int64(services.HashPassword(pass))
	email := string(user)

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		logs.ErrorLog(dbEngine.ErrDBNotFound)
		return nil
	}

	users, err := db.NewUsers(DB)
	if err != nil {
		logs.ErrorLog(err, "get table 'users' fail")
		return nil
	}
	u := &User{}
	err = users.SelectSelfScanEach(
		ctx,
		func(record *db.UsersFields) error {
			u.UsersFields = users.Record

			return nil
		},
		dbEngine.WhereForSelect("email", "hash"),
		dbEngine.ArgsForSelect(email, hash),
	)
	if err != nil {
		logs.ErrorLog(err, "select from table 'users' fail")
		return nil
	}
	if u.UsersFields == nil {
		logs.DebugLog("%s not found! (%d)", email, hash)
		return nil
	}
	return u

}
