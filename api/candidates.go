// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/ruslanBik4/uppeople/auth"
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
	MySent          bool          `json:"mySent"`
	Sort            int32         `json:"sort"`
	CurrentColumn   string        `json:"currentColumn"`
	DateFrom        string        `json:"dateFrom"`
	DateTo          string        `json:"dateTo"`
	SelectRecruiter *SelectedUnit `json:"selectRecruiter"`
	SelectCompanies SelectedUnits `json:"selectCompanies"`
	SelectReason    *SelectedUnit `json:"selectReason"`
	SelectTag       *SelectedUnit `json:"selectTag"`
	SelectPlatforms SelectedUnits `json:"selectPlatforms"`
	SelectSeniority SelectedUnits `json:"selectSeniority"`
	SelectStatuses  SelectedUnits `json:"selectStatuses"`
}

func (s *SearchCandidates) GetValue() interface{} {
	return s
}

func (s *SearchCandidates) NewValue() interface{} {
	return &SearchCandidates{}
}

func HandleAllCandidates(ctx *fasthttp.RequestCtx) (interface{}, error) {
	state := "start"

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	offset := 0
	id, ok := ctx.UserValue(ParamPageNum.Name).(int)
	if ok && id > 1 {
		offset = (id - 1) * pageItem
	}

	res := ResCandidates{
		ResList:    NewResList(ctx, DB, id),
		Candidates: make([]*CandidateView, 0),
		Company:    getCompanies(ctx, DB),
		Reasons:    getRejectReason(ctx, DB),
		Recruiter:  getRecruters(ctx, DB),
		Statuses:   getStatusVac(ctx, DB),
		Tags:       getTags(ctx, DB),
	}
	var err error
	ch := make(chan string, 0)
	go func() {

		defer func() {
			ch <- "end"
		}()
		optionsCount := []dbEngine.BuildSqlOptions{
			dbEngine.ColumnsForSelect("count(*)"),
		}
		options := []dbEngine.BuildSqlOptions{
			dbEngine.OrderBy("date desc"),
			dbEngine.FetchOnlyRows(pageItem),
			dbEngine.Offset(offset),
		}
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

			if dto.SelectTag != nil {
				if dto.SelectTag.Id == 3 {
					if dto.SelectReason != nil {
						where = append(where, "tag_id")
						args = append(args, dto.SelectReason.Id)
					} else {
						where = append(where, `tag_id in (SELECT id 
												FROM tags 
												WHERE parent_id=%s)`)
						args = append(args, dto.SelectTag.Id)
					}
				} else {
					where = append(where, "tag_id")
					args = append(args, dto.SelectTag.Id)
				}
			}

			if dto.SelectPlatforms != nil {
				where = append(where, "platform_id")
				p := make([]int32, len(dto.SelectPlatforms))
				for i, unit := range dto.SelectPlatforms {
					p[i] = unit.Id
				}
				args = append(args, p)
			}

			if dto.SelectSeniority != nil {
				where = append(where, "seniority_id")
				p := make([]int32, len(dto.SelectSeniority))
				for i, unit := range dto.SelectSeniority {
					p[i] = unit.Id
				}
				args = append(args, p)
			}

			if dto.DateFrom > "" {
				where = append(where, ">=date")
				args = append(args, dto.DateFrom)

			}
			if dto.DateTo > "" {
				where = append(where, "<=date")
				args = append(args, dto.DateTo+"T23:59:59+02:00")

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
				for i, unit := range dto.SelectStatuses {
					p[i] = unit.Id
				}
				args = append(args, p)
			}

			if len(dto.SelectCompanies) > 0 {
				where = append(where, `id in (SELECT vc.candidate_id 
	FROM vacancies v JOIN companies on (v.company_id=companies.id)
	JOIN vacancies_to_candidates vc on (v.id = vc.vacancy_id)
	WHERE companies.id=%s)`)
				p := make([]int32, len(dto.SelectCompanies))
				for i, unit := range dto.SelectCompanies {
					p[i] = unit.Id
				}
				args = append(args, p)
			}

			if _, ok := ctx.UserValue("sendCandidate").(bool); ok {
				where = append(where, `id in (SELECT candidate_id 
	FROM vacancies_to_candidates
	WHERE status != %s)`)
				args = append(args, 1)
				if dto.MySent {
					where = append(where, "recruter_id")
					args = append(args, auth.GetUserID(ctx))
				}
			}

			if dto.CurrentColumn > "" {
				orderBy := dto.CurrentColumn
				switch dto.CurrentColumn {
				case "Companies":
					orderBy = `(select name from vacancies_to_candidates v JOIN companies on (v.company_id=companies.id) 
where candidates.id = candidate_id
order by status desc
fetch first 1 row only)`
				case "Platform":
					orderBy = `(select nazva from platforms where id = platform_id)`
				case "Recruiter":
					orderBy = `(select name from users where id = recruter_id)`
				case "Tag/Reason":
					orderBy = `(select name from tags where id = tag_id)`
				case "Seniority":
					orderBy = `(select nazva from seniorities where id = seniority_id)`
				case "Status":
					orderBy = `(select status from vacancies_to_candidates where candidates.id = candidate_id
order by status desc
fetch first 1 row only
)`
				case "Contacts":
					orderBy = `coalesce(email, phone, skype)`
				}
				if dto.Sort > 0 {
					orderBy += " desc"
				}
				options = append(options, dbEngine.OrderBy(orderBy))
			}

		}

		if _, ok := ctx.UserValue("sendCandidate").(bool); ok {
			where = append(where, `id in (SELECT candidate_id 
	FROM vacancies_to_candidates
	WHERE status > %s)`)
			args = append(args, 1)
		}

		if len(where) > 0 {
			options = append(options,
				dbEngine.WhereForSelect(where...),
				dbEngine.ArgsForSelect(args...),
			)
			optionsCount = append(optionsCount,
				dbEngine.WhereForSelect(where...),
				dbEngine.ArgsForSelect(args...),
			)
		}

		state = "read candidates"
		candidates, _ := db.NewCandidates(DB)
		err = candidates.SelectSelfScanEach(ctx,
			func(record *db.CandidatesFields) error {
				state = "read " + record.Name
				cV := NewCandidateView(ctx, record, DB, res.Platforms, res.Seniority)
				res.Candidates = append(res.Candidates, cV)

				return nil
			},
			options...,
		)
		if err != nil {
			state = err.Error()
			return
		}

		if len(res.Candidates) == 0 {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		state = "count candidates"

		if id == 1 && len(res.Candidates) < pageItem {
			res.ResList.TotalPage = 1
			res.ResList.Count = len(res.Candidates)
		} else {
			err = candidates.SelectOneAndScan(ctx,
				&res.ResList.Count,
				optionsCount...)
			if err != nil {
				logs.ErrorLog(err, "count")
			} else {
				res.ResList.TotalPage = res.ResList.Count / pageItem
				if res.ResList.Count%pageItem > 0 {
					res.ResList.TotalPage++
				}
			}
		}
	}()

	select {
	case s, ok := <-ch:
		if ok {
			logs.DebugLog(s)
		}

	case <-time.After(time.Minute * 3):
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		logs.DebugLog(state)
		return nil, nil
	}

	if err != nil {
		return createErrResult(err)
	}
	if len(res.Candidates) == 0 {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return nil, nil
	}
	return res, nil
}

func HandleReturnAllCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	ctx.SetUserValue("sendCandidate", true)
	return HandleAllCandidates(ctx)
}

func HandleCandidatesFreelancerOnVacancies(ctx *fasthttp.RequestCtx) (interface{}, error) {
	ctx.SetStatusCode(fasthttp.StatusNoContent)
	return nil, nil
}

func HandleAllCandidatesForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	return HandleAllCandidates(ctx)
}
