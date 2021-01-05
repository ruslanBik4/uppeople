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

type vacDTO struct {
	SelectPlatforms       []SelectedUnit `json:"selectPlatforms"`
	SelectSeniorities     []SelectedUnit `json:"selectSeniorities"`
	SelectCandidateStatus []SelectedUnit `json:"selectCandidate_status"`
	SelectStatuses        []SelectedUnit `json:"selectStatuses"`
}

func (v *vacDTO) GetValue() interface{} {
	return v
}

func (v *vacDTO) NewValue() interface{} {
	return &vacDTO{}
}

type ResVacancies struct {
	*ResList
	CandidateStatus SelectedUnits         `json:"candidate_status"`
	VacancyStatus   SelectedUnits         `json:"vacancy_status"`
	Vacancies       []*db.VacanciesFields `json:"vacancies"`
}

func HandleViewAllVacancyInCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	filter, ok := ctx.UserValue(apis.JSONParams).(*vacDTO)
	if !ok {
		return "DTO is wrong", apis.ErrWrongParamsList
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)
	if l := len(filter.SelectStatuses); l > 0 {
		arg := make([]int32, l)
		for i, s := range filter.SelectStatuses {
			arg[i] = s.Id
		}
		columns = append(columns, "status")
		args = append(args, arg)
	}

	if l := len(filter.SelectPlatforms); l > 0 {
		arg := make([]int32, l)
		for i, s := range filter.SelectPlatforms {
			arg[i] = s.Id
		}
		columns = append(columns, "platform_id")
		args = append(args, arg)
	}

	// if l := len(filter.SelectCandidateStatus); l > 0 {
	// 	arg := make([]int32, l)
	// 	for i, s := range filter.SelectCandidateStatus {
	// 		arg[i] = s.Id
	// 	}
	// 	columns = append(columns, "status")
	// 	args = append(args, arg)
	// }
	//
	if l := len(filter.SelectSeniorities); l > 0 {
		arg := make([]int32, l)
		for i, s := range filter.SelectSeniorities {
			arg[i] = s.Id
		}
		columns = append(columns, "seniority_id")
		args = append(args, arg)
	}

	vacancies, _ := db.NewVacancies(DB)
	res := ResVacancies{
		ResList:         NewResList(ctx, DB),
		Vacancies:       make([]*db.VacanciesFields, 0),
		CandidateStatus: getStatusVac(ctx, DB),
		VacancyStatus:   getStatuses(ctx, DB),
	}

	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy("date_create desc"),
	}
	if len(columns) > 0 {
		options = append(options, dbEngine.WhereForSelect(columns...))
		options = append(options, dbEngine.ArgsForSelect(args...))
	}
	i := 0
	err := vacancies.SelectSelfScanEach(ctx,
		func(record *db.VacanciesFields) error {
			res.Vacancies = append(res.Vacancies, record)

			i++
			if i == pageItem {
				return errLimit
			}

			return nil
		},
		options...,
	)

	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	return res, nil

}
