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

func HandleDownloadWholeReport(ctx *fasthttp.RequestCtx) (interface{}, error) {
	sqlCmd := fmt.Sprintf("select * from whole_report('%s', '%s')",
		ctx.UserValue(ParamStartDate.Name),
		ctx.UserValue(ParamEndDate.Name),
	)
	return downloadCommand(ctx, sqlCmd)
}

func HandleDownloadReportByTag(ctx *fasthttp.RequestCtx) (interface{}, error) {
	s, err := createCommandWithSql(ctx, "amoung_by_tags")
	if err != nil {
		return createErrResult(err)
	}

	return s, nil
}

func HandleDownloadReportByStatus(ctx *fasthttp.RequestCtx) (interface{}, error) {
	s, err := createCommandWithSql(ctx, "amoung_by_status")
	if err != nil {
		return createErrResult(err)
	}

	return s, nil
}

func createCommandWithSql(ctx *fasthttp.RequestCtx, funcName string) (interface{}, error) {
	p, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong json", apis.ErrWrongParamsList
	}

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
	sqlCmd := fmt.Sprintf(`select id, name, "count", percent
							from %s('%s', '%s', %d, %d, %d, %d, %s)`,
		funcName,
		p.StartDate,
		p.EndDate,
		p.RecruiterId,
		p.CompanyId,
		p.PlatformId,
		p.VacancyId,
		arr)

	return downloadCommand(ctx, sqlCmd)
}

func downloadCommand(ctx *fasthttp.RequestCtx, sqlCmd string) (interface{}, error) {
	sqlCmd = fmt.Sprintf(`\copy ( %s )
						to stdout csv header;`, sqlCmd)
	// todo add DB name from connection
	cmd := exec.CommandContext(ctx, `sudo`, `-u`, `postgres`, `psql`, "-d", "uppeople", "-c", sqlCmd)
	return cmdResponseAsFile(ctx, cmd)
}

func cmdResponseAsFile(ctx *fasthttp.RequestCtx, cmd *exec.Cmd) (interface{}, error) {
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
