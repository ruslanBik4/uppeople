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
	CurrentPage            int              `json:"currentPage"`
	PerPage                int              `json:"perPage"`
	Platforms              db.SelectedUnits `json:"platforms"`
	Seniority              db.SelectedUnits `json:"seniority"`
}

func NewCandidateView(ctx *fasthttp.RequestCtx,
	record *db.CandidatesFields,
	DB *dbEngine.DB,
	platforms db.SelectedUnits,
	seniors db.SelectedUnits,
) *CandidateView {
	ref := &CandidateView{
		ViewCandidate: &ViewCandidate{
			CandidatesFields: record,
		},
		Status: statusCandidate{
			Date:         record.Date,
			Comments:     record.Comments,
			DateFollowUp: record.Date_follow_up,
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
	} else {
		ref.Status.Recruiter = ref.Recruiter
	}

	ref.Tags = db.GetTagFromId(record.Tag_id)
	if ref.Tags != nil {
		ref.TagName = ref.Tags.Name
		ref.TagColor = ref.Tags.Color
	}

	ref.Seniority = db.GetSeniorityFromId(record.Seniority_id).Nazva.String

	for _, p := range platforms {
		if p.Id == record.Platform_id {
			ref.Platform = p.Label
			ref.ViewCandidate.Platform = p
		}
	}

	ref.ViewCandidate.Vacancies, err = DB.Conn.SelectToMaps(ctx,
		`select v.id, 
		concat(companies.name, ' ("', platforms.nazva, '")') as name, 
		concat(companies.name, ' ("', platforms.nazva, '")') as label, 
		LOWER(CONCAT(companies.name, ' ("', platforms.nazva , '")')) as value, 
		user_ids, 
		platform_id,
		CONCAT(platforms.nazva, ' ("', (select nazva from seniorities where id=seniority_id), '")') as platform,
		companies, sv.id as status_id, v.company_id, sv.status, salary, 
		coalesce(vc.date_last_change, vc.date_create) as date_last_change, vc.rej_text, sv.color
FROM vacancies v JOIN companies on (v.company_id=companies.id)
	JOIN vacancies_to_candidates vc on (v.id = vc.vacancy_id )
	JOIN platforms ON (v.platform_id = platforms.id)
	JOIN status_for_vacs sv on (vc.status = sv.id)
	WHERE vc.candidate_id=$1 AND vc.status!=1
    order by date_last_change desc
`,
		ref.ViewCandidate.Id,
	)
	if err != nil {
		logs.ErrorLog(err, "")
	}

	if len(ref.ViewCandidate.Vacancies) > 0 {
		vacancy := ref.ViewCandidate.Vacancies[0]
		ref.Status.CompId, _ = vacancy["company_id"].(int32)
		ref.Status.CompName, _ = vacancy["name"].(string)
		ref.Status.Comments, _ = vacancy["rej_text"].(string)
		ref.Status.VacStat, _ = vacancy["status"].(string)
		ref.Status.Date, _ = vacancy["date_last_change"].(time.Time)
		ref.Color, _ = vacancy["color"].(string)

		args := make([]int32, len(ref.ViewCandidate.Vacancies))
		for i, vac := range ref.ViewCandidate.Vacancies {
			args[i] = vac["company_id"].(int32)
		}
		ref.Companies = getCompanies(ctx,
			DB,
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(args),
		)
	}

	return ref
}

func NewResList(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, pageNum int) *ResList {
	return &ResList{
		Page:        pageNum,
		CurrentPage: pageNum,
		Count:       10 * pageItem,
		TotalPage:   10,
		PerPage:     pageItem,
		Platforms:   getPlatforms(ctx, DB),
		Seniority:   db.GetSenioritiesAsSelectedUnits(),
	}
}
