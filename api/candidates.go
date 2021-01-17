// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/pkg/errors"
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
	Name     string `json:"search"`
	DateFrom string `json:"date_from"`
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

	seniors := getSeniorities(ctx, DB)
	args := make([]interface{}, 0)
	where := make([]string, 0)

	dto, ok := ctx.UserValue(apis.JSONParams).(*SearchCandidates)
	if ok {
		if dto.Name > "" {
			where = append(where, "~name")
			args = append(args, dto.Name)
		}
		if dto.DateFrom > "" {
			where = append(where, "date")
			args = append(args, dto.DateFrom)

		}
	}
	offset := 0
	id, ok := ctx.UserValue("id").(int32)
	if ok && id > 1 {
		offset = int(id * pageItem)
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
		return nil, errors.Wrap(err, "	")
	}

	return res, nil
}
