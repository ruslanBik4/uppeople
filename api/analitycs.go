// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
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
	Recruiter_id int32  `json:"recruiter_id"`
	Company_id   int32  `json:"company_id"`
	Vacancy_id   int32  `json:"vacancy_id"`
	Start_date   string `json:"start_date"`
	End_date     string `json:"end_date"`
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
	Main, Reject []map[string]interface{}
}

func HandleGetCandidatesAmountByTags(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	sql := `SELECT id, name, color, count(c.id), parent_id
	FROM tags t
			LEFT JOIN candidates c ON t.id=c.tag_id
			LEFT JOIN  vacancies_to_candidates vtc ON c.id=vtc.candidate_id
			LEFT JOIN vacancies v ON v.id = vtc.vacancy_id
			WHERE v.status IN (0,1)
       GROUP BY 1, 2, 3, 5`
	m, err := DB.Conn.SelectToMaps(ctx,
		sql,
	)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	r, err := DB.Conn.SelectToMaps(ctx,
		sql,
	)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return AmountsByTags{
		Main:   m,
		Reject: r,
	}, nil
}
