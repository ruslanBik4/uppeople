// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type ResCandidates struct {
	*ResList
	Candidates []*CandidateView `json:"candidates"`
	Company    SelectedUnits    `json:"company"`
	Recruiter  SelectedUnits    `json:"recruiter"`
	Reasons    SelectedUnits    `json:"reasons"`
	Statuses   SelectedUnits    `json:"statuses"`
	Tags       SelectedUnits    `json:"tags"`
}

type SearchCandidates struct {
	Name            string        `json:"search"`
	DateFrom        string        `json:"dateFrom"`
	DateTo          string        `json:"dateTo"`
	SelectRecruiter *SelectedUnit `json:"selectRecruiter"`
	SelectCompanies SelectedUnits `json:"selectCompanies"`
}

func (s *SearchCandidates) GetValue() interface{} {
	return s
}

func (s *SearchCandidates) NewValue() interface{} {
	return &SearchCandidates{}
}

func HandleAllCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	candidates, _ := db.NewCandidates(DB)
	res := ResCandidates{
		ResList:    NewResList(ctx, DB),
		Candidates: make([]*CandidateView, 0),
		Company:    getCompanies(ctx, DB),
		Reasons:    getRejectReason(ctx, DB),
		Recruiter:  getRecruters(ctx, DB),
		Statuses:   getStatuses(ctx, DB),
		Tags:       getTags(ctx, DB),
	}

	offset := 0
	id, ok := ctx.UserValue("id").(int32)
	if ok && id > 1 {
		offset = int(id * pageItem)
	}

	seniors := getSeniorities(ctx, DB)
	args := make([]interface{}, 0)
	where := make([]string, 0)

	dto, ok := ctx.UserValue(apis.JSONParams).(*SearchCandidates)
	if ok {
		if dto.Name > "" {
			where = append(where, "~name")
			args = append(args, dto.Name)
		}
		if dto.SelectRecruiter != nil {
			where = append(where, "recruter_id")
			args = append(args, dto.SelectRecruiter.Id)
		}
		if dto.DateFrom > "" {
			where = append(where, ">=date")
			args = append(args, dto.DateFrom)

		}
		if dto.DateTo > "" {
			where = append(where, "<=date")
			args = append(args, dto.DateTo)

		}
		if len(dto.SelectCompanies) > 0 {
			where = append(where, "companies.id")
			args = append(args, dto.SelectRecruiter.Id)
			// DB.Conn.SelectToMaps(ctx,
			// 	`select candidates.id',
			//         GROUP_CONCAT(DISTINCT(JSON_OBJECT(compId, companies.id, compName, companies.nazva,  vacStat, status_for_vacs.status, date, IFNULL(vacancies_to_candidates.date_last_change, vacancies_to_candidates.date_create), commentVac, vacancies_to_candidates.rej_text)) SEPARATOR ;) as status'),
			//         candidates.date,
			//         platforms.nazva as platform,
			//         candidates.name,
			//         candidates.email,
			//         candidates.mobile,
			//         candidates.skype,
			//         candidates.linkedin,
			//         salary,
			//         status_for_vacs.color,
			//         GROUP_CONCAT(DISTINCT(JSON_OBJECT(id, companies.id, name, companies.nazva, vacStat, status_for_vacs.status)) SEPARATOR ;) as companies'),
			//         users.name as recruiter`,
			// 	args...)
		}
	}
	err := candidates.SelectSelfScanEach(ctx,
		func(record *db.CandidatesFields) error {

			cV := NewCandidateView(ctx, record, DB, res.Platforms, seniors)
			res.Candidates = append(res.Candidates, cV)

			return nil
		},
		dbEngine.WhereForSelect(where...),
		dbEngine.ArgsForSelect(args...),
		dbEngine.OrderBy("date desc"),
		dbEngine.ArgsForSelect(args...),
		dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	)
	if err != nil {
		return createErrResult(err)
	}

	return res, nil
}

func HandleReturnAllCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	return HandleAllCandidate(ctx)
}
