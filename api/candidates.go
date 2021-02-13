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
	CompanyID       int32         `json:"company_id"`
	DateFrom        string        `json:"dateFrom"`
	DateTo          string        `json:"dateTo"`
	SelectRecruiter *SelectedUnit `json:"selectRecruiter"`
	SelectCompanies SelectedUnits `json:"selectCompanies"`
	SelectTag       *SelectedUnit `json:"selectTag"`
	SelectPlatforms SelectedUnits `json:"selectPlatforms"`
	SelectStatuses  SelectedUnits `json:"selectStatuses"`
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

	offset := 0
	id, ok := ctx.UserValue(ParamPageNum.Name).(int)
	if ok && id > 1 {
		offset = id * pageItem
	}

	candidates, _ := db.NewCandidates(DB)
	res := ResCandidates{
		ResList:    NewResList(ctx, DB, id),
		Candidates: make([]*CandidateView, 0),
		Company:    getCompanies(ctx, DB),
		Reasons:    getRejectReason(ctx, DB),
		Recruiter:  getRecruters(ctx, DB),
		Statuses:   getStatusVac(ctx, DB),
		Tags:       getTags(ctx, DB),
	}

	seniors := getSeniorities(ctx, DB)
	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy("date desc"),
		dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	}
	dto, ok := ctx.UserValue(apis.JSONParams).(*SearchCandidates)
	if ok {
		args := make([]interface{}, 0)
		where := make([]string, 0)

		if dto.Name > "" {
			where = append(where, "~name")
			args = append(args, dto.Name)
		}
		if dto.SelectRecruiter != nil {
			where = append(where, "recruter_id")
			args = append(args, dto.SelectRecruiter.Id)
		}

		if dto.SelectTag != nil {
			where = append(where, "tag_id")
			args = append(args, dto.SelectTag.Id)
		}

		if dto.SelectPlatforms != nil {
			where = append(where, "platform_id")
			p := make([]int32, len(dto.SelectPlatforms))
			for i, tag := range dto.SelectPlatforms {
				p[i] = tag.Id
			}
			args = append(args, p)
		}

		if dto.DateFrom > "" {
			where = append(where, ">=date")
			args = append(args, dto.DateFrom)

		}
		if dto.DateTo > "" {
			where = append(where, "<=date")
			args = append(args, dto.DateTo)

		}
		if dto.CompanyID > 0 {
			where = append(where, `id in (SELECT candidate_id 
					FROM candidates_to_companies
					WHERE  company_id = %s)	`)
			args = append(args, dto.CompanyID)
		}
		if dto.SelectStatuses != nil {
			where = append(where, `id in (SELECT candidate_id 
	FROM vacancies_to_candidates
	WHERE status=%s)`)
			p := make([]int32, len(dto.SelectStatuses))
			for i, tag := range dto.SelectStatuses {
				p[i] = tag.Id
			}
			args = append(args, p)
		}

		if len(dto.SelectCompanies) > 0 {
			where = append(where, `id in (SELECT vc.candidate_id 
	FROM vacancies v JOIN companies on (v.company_id=companies.id)
	JOIN vacancies_to_candidates vc on (v.id = vc.vacancy_id)
	WHERE companies.id=%s)`)
			p := make([]int32, len(dto.SelectCompanies))
			for i, tag := range dto.SelectCompanies {
				p[i] = tag.Id
			}
			args = append(args, p)
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
		options = append(options,
			dbEngine.WhereForSelect(where...),
			dbEngine.ArgsForSelect(args...),
		)

	}
	err := candidates.SelectSelfScanEach(ctx,
		func(record *db.CandidatesFields) error {

			cV := NewCandidateView(ctx, record, DB, res.Platforms, seniors)
			res.Candidates = append(res.Candidates, cV)

			return nil
		},
		options...,
	)
	if err != nil {
		return createErrResult(err)
	}

	if len(res.Candidates) < pageItem {
		res.ResList.TotalPage = 1
		res.ResList.Count = len(res.Candidates)
	}

	return res, nil
}

func HandleReturnAllCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	return HandleAllCandidate(ctx)
}

func HandleAllCandidatesForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	return HandleAllCandidate(ctx)
}
