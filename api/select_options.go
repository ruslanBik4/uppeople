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

	res := NewSelectOpt(ctx, DB)

	return res, nil
}

type selectOpt struct {
	Companies     SelectedUnits `json:"companies"`
	Platforms     SelectedUnits `json:"platforms"`
	Recruiters    SelectedUnits `json:"recruiters"`
	Statuses      SelectedUnits `json:"candidateStatus"`
	Location      SelectedUnits `json:"location"`
	RejectReasons SelectedUnits `json:"reject_reasons"`
	RejectTag     SelectedUnits `json:"reject_tag"`
	Recruiter     SelectedUnits `json:"recruiter"`
	Seniorities   SelectedUnits `json:"seniorities"`
	Tags          SelectedUnits `json:"tags"`
	VacancyStatus SelectedUnits `json:"vacancyStatus"`
	// vacancies
}

func NewSelectOpt(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) selectOpt {
	s := selectOpt{
		Companies:     getCompanies(ctx, DB),
		Platforms:     getPlatforms(ctx, DB),
		Recruiters:    getRecruters(ctx, DB),
		Statuses:      getStatusVac(ctx, DB),
		Location:      getLocations(ctx, DB),
		RejectReasons: getRejectReason(ctx, DB),
		Seniorities:   getSeniorities(ctx, DB),
		Tags:          getTags(ctx, DB),
		VacancyStatus: getStatusVac(ctx, DB),
	}
	for _, tag := range s.Tags {
		if tag.Id == 3 {
			s.RejectTag = SelectedUnits{tag}
		}
	}

	return s
}
