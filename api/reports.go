// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

func HandleDownloadExecutions(ctx *fasthttp.RequestCtx) (interface{}, error) {
	p, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong json", apis.ErrWrongParamsList
	}

	sqlCmd := fmt.Sprintf(`\copy (
    select *
    from amoung_by_tags('%s', '%s', %d, %d, %d)
)
to stdout csv header;`, p.StartDate, p.EndDate, p.RecruiterId, p.CompanyId, p.VacancyId)
	cmd := exec.CommandContext(ctx, `sudo`, `-u`, `postgres`, `psql`, "-d", "test", "-c", sqlCmd)
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
