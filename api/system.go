// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"os/exec"
	"runtime/trace"
	"strings"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

func HandleTrace(ctx *fasthttp.RequestCtx) (interface{}, error) {
	trace.Stop()
	cmd := exec.CommandContext(ctx, `sudo`, `go`, `tool`, `trace`, "dev.out")
	cmd.Stdout = ctx.Response.BodyWriter()
	cmd.Stderr = ctx.Response.BodyWriter()
	err := cmd.Run()
	s := string(ctx.Response.Body())
	if err != nil {
		logs.DebugLog(s)

		return s, errors.Wrap(err, cmd.String()+s)
	}

	parts := strings.Split(s, ":")
	url := fmt.Sprintf("%s://%s%s/", ctx.URI().Scheme(), string(ctx.Host()), parts[len(parts)-1])
	logs.DebugLog("redirect:", url)
	ctx.Redirect(url, fasthttp.StatusMovedPermanently)

	return nil, nil
}

func HandleStatConn(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return DB.Conn.GetStat(), nil
}
