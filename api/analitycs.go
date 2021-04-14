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
}

func (d *DTOAmounts) GetValue() interface{} {
	return d
}

func (d *DTOAmounts) NewValue() interface{} {
	return &DTOAmounts{}
}

type Amoint struct {
}
type AmountsByTags struct {
	Message string                   `json:"message"`
	Data    interface{}              `json:"data,omitempty"`
	Main    []map[string]interface{} `json:"main,omitempty"`
	Reject  []map[string]interface{} `json:"reject,omitempty"`
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

	sql := `SELECT vtc.status as status_id, sfv.status, sfv.color, count(vtc.id) as count 
			FROM vacancies_to_candidates vtc
			 JOIN status_for_vacs sfv ON sfv.id = vtc.status
			 JOIN candidates c ON c.id = vtc.candidate_id
			 JOIN vacancies v ON v.id = vtc.vacancy_id
			WHERE v.status IN (0,1) AND vtc.status > 1
`
	gr := `      GROUP BY 1,2,3`

	params, ok := ctx.UserValue(apis.JSONParams).(*DTOAmounts)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	where := getParams(params)

	m, err := DB.Conn.SelectToMaps(ctx,
		sql+where+gr,
	)
	if err != nil {
		return createErrResult(err)
	}

	data := make(map[string]map[string]interface{}, len(m))
	for _, val := range m {
		data[val["status"].(string)] = val
	}
	return AmountsByTags{
		Message: "Successfully",
		Data:    data,
	}, nil
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

	m := make([]map[string]interface{}, 0)
	r := make([]map[string]interface{}, 0)
	err := proc.SelectAndRunEach(ctx,
		func(values []interface{}, columns []dbEngine.Column) error {
			row := make(map[string]interface{})
			for i, col := range columns {
				row[col.Name()] = values[i]
			}

			if row["parent_id"] == 0 {
				m = append(m, row)
			} else {
				r = append(r, row)
			}
			return nil
		},
		dbEngine.ArgsForSelect(params.StartDate, params.EndDate, params.RecruiterId, params.CompanyId, params.VacancyId),
	)
	if err != nil {
		return createErrResult(err)
	}

	return AmountsByTags{
		Message: "Successfully",
		Main:    m,
		Reject:  r,
	}, nil
}

func getParams(params *DTOAmounts) string {
	where := ""
	if p := params.CompanyId; p > 0 {
		where += fmt.Sprintf(" and v.company_id = %d", p)
	}

	if p := params.VacancyId; p > 0 {
		where += fmt.Sprintf(" and v.id = %d", p)
	}

	if p := params.RecruiterId; p > 0 {
		where += fmt.Sprintf(" and c.recruter_id = %d", p)
	}

	if p := params.StartDate; p > "" {
		where += fmt.Sprintf(" AND c.date >= '%s'", p)
	}

	if p := params.EndDate; p > "" {
		where += fmt.Sprintf(" AND c.date <= '%s'", p)
	}

	return where
}
