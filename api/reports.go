// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

func HandleDownloadReportByTag(ctx *fasthttp.RequestCtx) (interface{}, error) {
	p, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong json", apis.ErrWrongParamsList
	}

	s, err := createCommandWithSql(ctx, "amoung_by_tags", p)
	if err != nil {
		return createErrResult(err)
	}

	return s, nil
}

func HandleDownloadReportByStatus(ctx *fasthttp.RequestCtx) (interface{}, error) {
	p, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong json", apis.ErrWrongParamsList
	}

	s, err := createCommandWithSql(ctx, "amoung_by_status", p)
	if err != nil {
		return createErrResult(err)
	}

	return s, nil
}

func createCommandWithSql(ctx *fasthttp.RequestCtx, funcName string, p *DTOAmounts) (interface{}, error) {
	arr := "NULL"
	if len(p.Includes) > 0 {
		arr = "ARRAY["
		for key, val := range p.Includes {
			if key > 0 {
				arr += ","
			}
			arr += strconv.Itoa(val)
		}
		arr += "]"
	}
	sqlCmd := fmt.Sprintf(`\copy (
							select *
							from %s('%s', '%s', %d, %d, %d, %s)
						)
						to stdout csv header;`,
		funcName,
		p.StartDate,
		p.EndDate,
		p.RecruiterId,
		p.CompanyId,
		p.VacancyId,
		arr)

	return downloadCommand(ctx, sqlCmd)
}

func downloadCommand(ctx *fasthttp.RequestCtx, sqlCmd string) (interface{}, error) {
	// todo add DB name from connection
	cmd := exec.CommandContext(ctx, `sudo`, `-u`, `postgres`, `psql`, "-d", "postgres", "-c", sqlCmd)
	cmd.Stdout = ctx.Response.BodyWriter()
	cmd.Stderr = ctx.Response.BodyWriter()
	err := cmd.Run()
	if err != nil {
		s := string(ctx.Response.Body())
		logs.DebugLog(s)

		return s, errors.Wrap(err, cmd.String()+s)
	}

	ctx.Response.Header.Set("Content-Description", "File Transfer")
	ctx.SetContentType("application/octet-stream")
	ctx.Response.Header.Set("Content-Transfer-Encoding", "binary")
	ctx.Response.Header.Set("Cache-Control", "must-revalidate")
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename=anton.csv")

	ctx.SetStatusCode(fasthttp.StatusOK)

	return nil, nil
}
