// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"
)

func HandleGetTags(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return getTags(ctx, DB), nil
}

func HandleGetStatuses(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return getStatuses(ctx, DB), nil
}

type DTOAmounts struct {
	RecruiterId int32  `json:"recruiter_id"`
	CompanyId   int32  `json:"company_id"`
	VacancyId   int32  `json:"vacancy_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Includes    []int  `json:"includes"`
}

func (d *DTOAmounts) GetValue() interface{} {
	return d
}

func (d *DTOAmounts) NewValue() interface{} {
	return &DTOAmounts{}
}

func (d *DTOAmounts) GetParamsArgs() dbEngine.BuildSqlOptions {
	return dbEngine.ArgsForSelect(
		d.StartDate,
		d.EndDate,
		d.RecruiterId,
		d.CompanyId,
		d.VacancyId,
		d.Includes,
	)
}

type AmountsByTags struct {
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data,omitempty"`
	Main    []map[string]interface{} `json:"main,omitempty"`
	Reject  []map[string]interface{} `json:"reject,omitempty"`
	Total   int32                    `json:"total"`
}

func HandleGetCandidatesByVacancies(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	sql := `SELECT vacancy_id, vtc.company_id,
				 p.nazva as platform_name,
				  c.name as company_name,
				   u.name as user_name,
				   user_id as recruiter_id, 
				count(candidate_id) as quantity
		FROM vacancies_to_candidates vtc
				LEFT JOIN vacancies v ON v.id=vtc.vacancy_id
				LEFT JOIN platforms p ON p.id=v.platform_id
				LEFT JOIN companies c ON c.id=vtc.company_id
				LEFT JOIN users u ON u.id=vtc.user_id
				WHERE v.status IN (0,1) `
	gr := ` GROUP BY 1,2,3,4,5,6 ORDER BY 2`

	params, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	where := ""
	if p := params.CompanyId; p > 0 {
		where += fmt.Sprintf(" and vtc.company_id = %d", p)
	}

	if p := params.RecruiterId; p > 0 {
		where += fmt.Sprintf(" and vtc.user_id = %d", p)
	}

	data, err := DB.Conn.SelectToMaps(ctx,
		sql+where+gr,
	)
	if err != nil {
		return createErrResult(err)
	}

	return AmountsByTags{
		Message: "Successfully",
		Data:    data,
	}, nil

}

func HandleGetCandidatesAmountByStatuses(ctx *fasthttp.RequestCtx) (interface{}, error) {

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	proc, ok := DB.Routines["amoung_by_status"]
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	params, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	res := AmountsByTags{
		Message: "Successfully",
		Data:    make([]map[string]interface{}, 0),
	}
	err := proc.SelectAndRunEach(ctx,
		func(values []interface{}, columns []dbEngine.Column) error {
			row := make(map[string]interface{})
			for i, col := range columns {
				row[col.Name()] = values[i]
			}

			if row["id"] != nil {
				res.Data = append(res.Data, row)
			} else {
				res.Total = row["count"].(int32)
			}

			return nil
		},
		params.GetParamsArgs(),
	)
	if err != nil {
		return createErrResult(err)
	}

	return res, nil
}

func HandleGetCandidatesAmountByTags(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	proc, ok := DB.Routines["amoung_by_tags"]
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	params, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	res := AmountsByTags{
		Message: "Successfully",
		Main:    make([]map[string]interface{}, 0),
		Reject:  make([]map[string]interface{}, 0),
	}
	err := proc.SelectAndRunEach(ctx,
		func(values []interface{}, columns []dbEngine.Column) error {
			row := make(map[string]interface{})
			for i, col := range columns {
				row[col.Name()] = values[i]
			}

			switch row["parent_id"] {
			case 3:
				res.Reject = append(res.Reject, row)
			case nil:
				res.Total = row["count"].(int32)
			default:
				res.Main = append(res.Main, row)
			}
			return nil
		},
		params.GetParamsArgs(),
	)
	if err != nil {
		return createErrResult(err)
	}

	return res, nil
}
