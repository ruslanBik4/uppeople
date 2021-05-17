// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type ResList struct {
	Count, Page, TotalPage int
	CurrentPage            int `json:"currentPage"`
	PerPage                int `json:"perPage"`
	// Platforms              db.SelectedUnits `json:"platforms"`
	// Seniority              db.SelectedUnits `json:"seniority"`
}

func NewResList(pageNum int) *ResList {
	return &ResList{
		Page:        pageNum,
		CurrentPage: pageNum,
		Count:       10 * pageItem,
		TotalPage:   10,
		PerPage:     pageItem,
		// Platforms:   db.GetPlatformsAsSelectedUnits(),
		// Seniority:   db.GetSenioritiesAsSelectedUnits(),
	}
}

func NewCandidateView(ctx *fasthttp.RequestCtx,
	record *db.CandidatesFields,
	DB *dbEngine.DB,
) *CandidateView {
	ref := &CandidateView{
		ViewCandidate: &ViewCandidate{
			CandidatesFields: record,
		},
	}

	users, _ := db.NewUsers(DB)

	err := users.SelectOneAndScan(ctx,
		&ref.Recruiter,
		dbEngine.ColumnsForSelect("name"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(record.Recruter_id),
	)
	if err != nil && err != pgx.ErrNoRows {
		logs.ErrorLog(err, "users.SelectOneAndScan")
	}

	ref.Tags = db.GetTagFromId(record.Tag_id)
	if ref.Tags != nil {
		ref.TagName = ref.Tags.Name
		ref.TagColor = ref.Tags.Color
	}

	ref.Seniority = db.GetSeniorityFromId(record.Seniority_id).Name

	platform := db.GetPlatformFromId(record.Platform_id)
	ref.Platform = platform.Name
	ref.ViewCandidate.Platform = &db.SelectedUnit{
		Id:    platform.Id,
		Label: platform.Name,
		Value: strings.ToLower(platform.Name),
	}

	ref.ViewCandidate.Vacancies, err = DB.Conn.SelectToMaps(ctx,
		SQL_VIEW_CANDIDATE_VACANCIES,
		ref.ViewCandidate.Id,
	)
	if err != nil {
		logs.ErrorLog(err, "")
	}

	args := make([]int32, len(ref.ViewCandidate.Vacancies))
	for i, vacancy := range ref.ViewCandidate.Vacancies {
		args[i] = vacancy["company_id"].(int32)
		ref.Statuses = append(ref.Statuses,
			&statusCandidate{
				CompId:    vacancy["company_id"].(int32),
				CompName:  vacancy["name"].(string),
				Comments:  vacancy["rej_text"].(string),
				VacStat:   vacancy["status"].(string),
				Date:      vacancy["date_last_change"].(time.Time),
				Color:     vacancy["color"].(string),
				Recruiter: vacancy["recruiter"].(string),
			})
	}
	if len(args) > 0 {
		ref.Companies = getCompanies(ctx,
			DB,
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(args),
		)
	}

	return ref
}
