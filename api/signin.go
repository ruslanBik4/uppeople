// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
)

func HandleAuthLogin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	req := &fasthttp.Request{}
	ctx.Request.CopyTo(req)

	for _, host := range hosts {
		resp := &fasthttp.Response{}

		err := doRequest(req, resp, host)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode() == fasthttp.StatusOK {
			var v map[string]interface{}
			err := jsoniter.Unmarshal(resp.Body(), &v)
			u := &auth.User{
				Companies: make(map[int32]map[string]string),
				Host:      host,
				TokenOld:  v["access_token"].(string),
				// UsersFields: users.Record,
			}

			u.Token, err = auth.Bearer.NewToken(u)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusNonAuthoritativeInfo)
				return nil, errors.Wrap(err, "Bearer.NewToken")
			}

			v["access_token"] = u.Token

			return v, nil
		}
	}

	return nil, nil
}
