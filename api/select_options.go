// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
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
	Companies     db.SelectedUnits `json:"companies"`
	Platforms     db.SelectedUnits `json:"platforms"`
	Recruiters    db.SelectedUnits `json:"recruiters"`
	Statuses      db.SelectedUnits `json:"candidateStatus"`
	Location      db.SelectedUnits `json:"location"`
	RejectReasons db.SelectedUnits `json:"reject_reasons"`
	RejectTag     db.SelectedUnits `json:"reject_tag"`
	Recruiter     db.SelectedUnits `json:"recruiter"`
	Seniorities   db.SelectedUnits `json:"seniorities"`
	Tags          db.SelectedUnits `json:"tags"`
	VacancyStatus db.SelectedUnits `json:"vacancyStatus"`
	// vacancies
}

func NewSelectOpt(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) selectOpt {
	s := selectOpt{
		Companies:     getCompanies(ctx, DB),
		Platforms:     getPlatforms(ctx, DB),
		Recruiters:    getRecruiters(ctx, DB),
		Statuses:      db.GetStatusForVacAsSelectedUnits(),
		Location:      getLocations(ctx, DB),
		RejectReasons: db.GetRejectReasonAsSelectedUnits(),
		Seniorities:   db.GetSenioritiesAsSelectedUnits(),
		Tags:          db.GetTagsAsSelectedUnits(),
		VacancyStatus: db.GetStatusAsSelectedUnits(),
	}
	for _, tag := range s.Tags {
		if tag.Id == db.GetTagIdReject() {
			s.RejectTag = db.SelectedUnits{tag}
		}
	}

	return s
}
