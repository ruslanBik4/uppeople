// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"os/exec"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

func HandleTrace(ctx *fasthttp.RequestCtx) (interface{}, error) {
	cmd := exec.CommandContext(ctx, `sudo`, `go`, `tool`, `trace`, "dev.out")
	cmd.Stdout = ctx.Response.BodyWriter()
	cmd.Stderr = ctx.Response.BodyWriter()
	err := cmd.Run()
	if err != nil {
		s := string(ctx.Response.Body())
		logs.DebugLog(s)

		return s, errors.Wrap(err, cmd.String()+s)
	}

	return nil, nil
}

func HandleStatConn(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return DB.Conn.GetStat(), nil
}
