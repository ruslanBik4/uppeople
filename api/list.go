// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
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
	view := &CandidateView{
		ViewCandidate: &ViewCandidate{
			CandidatesFields: record,
		},
	}

	users, _ := db.NewUsers(DB)

	err := users.SelectOneAndScan(ctx,
		&view.Recruiter,
		dbEngine.ColumnsForSelect("name"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(record.RecruterId),
	)
	if err != nil && err != pgx.ErrNoRows {
		logs.ErrorLog(err, "users.SelectOneAndScan")
	}

	selectedVacancies, ok := DB.Tables["select_vacancies"]
	if ok {
		err = selectedVacancies.SelectAndScanEach(ctx,
			nil,
			&view.SelectedVacancies,
			dbEngine.ColumnsForSelect("id", "label", "value"),
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(view.CandidatesFields.Vacancies),
		)
		if err != nil {
			logs.ErrorLog(err, "SelectedVacancies")
		}
	}

	view.Tags = db.GetTagFromId(record.TagId)
	if view.Tags != nil {
		view.TagName = view.Tags.Name
		view.TagColor = view.Tags.Color
	}

	view.Seniority = db.GetSeniorityFromId(record.SeniorityId).Name

	view.ViewCandidate.Vacancies, err = DB.Conn.SelectToMaps(ctx,
		SQL_VIEW_CANDIDATE_VACANCIES,
		view.ViewCandidate.Id,
	)
	if err != nil {
		logs.ErrorLog(err, "")
	}

	args := make([]int32, len(view.ViewCandidate.Vacancies))
	for i, vacancy := range view.ViewCandidate.Vacancies {
		args[i] = vacancy["company_id"].(int32)
		view.Statuses = append(view.Statuses,
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
		view.Companies = getCompanies(ctx,
			DB,
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(args),
		)
	}

	return view
}
