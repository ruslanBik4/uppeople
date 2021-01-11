// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/valyala/fasthttp"
)

func HandleReturnOptionsForSelects(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	res := selectOpt{
		Companies:     getCompanies(ctx, DB),
		Platforms:     getPlatforms(ctx, DB),
		Recruiters:    getRecruters(ctx, DB),
		Statuses:      getStatuses(ctx, DB),
		Location:      getLocations(ctx, DB),
		RejectReasons: getRejectReason(ctx, DB),
		Seniorities:   getSeniorities(ctx, DB),
		Tags:          getTags(ctx, DB),
		VacancyStatus: getStatusVac(ctx, DB),
	}

	for _, tag := range res.Tags {
		if tag.Id == 3 {
			res.RejectTag = SelectedUnits{tag}
		}
	}

	return res, nil
}
