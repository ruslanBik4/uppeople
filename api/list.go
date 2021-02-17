// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

func NewCandidateView(ctx *fasthttp.RequestCtx,
	record *db.CandidatesFields,
	DB *dbEngine.DB,
	platforms SelectedUnits,
	seniors SelectedUnits,
) *CandidateView {
	ref := &CandidateView{
		ViewCandidate: &ViewCandidate{
			CandidatesFields: record,
			Tags:             &db.TagsFields{},
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
	if err != nil {
		logs.ErrorLog(err, "users.SelectOneAndScan")
	}
	ref.Status.Recruiter = ref.Recruiter

	tagTable, _ := db.NewTags(DB)

	err = tagTable.SelectOneAndScan(ctx,
		ref.Tags,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(record.Tag_id),
	)
	if err != nil {
		logs.ErrorLog(err, "tagTable id=%d %s", record.Tag_id, record.Name)
	} else {
		ref.TagName = ref.Tags.Name
		ref.TagColor = ref.Tags.Color
	}

	for _, s := range seniors {
		if s.Id == ref.Seniority_id {
			ref.Seniority = s.Label
		}
	}

	for _, p := range platforms {
		if p.Id == record.Platform_id {
			ref.Platform = p.Label
			ref.ViewCandidate.Platform = p
		}
	}

	ref.ViewCandidate.Vacancies, err = DB.Conn.SelectToMaps(ctx,
		`select vacancies.id, concat(companies.name, ' ("', platforms.nazva, '")') as name, 
LOWER(CONCAT(companies.name, ' ("', platforms.nazva , ')"')) as label, user_ids, platform_id,
		companies, sv.id as status_id, vacancies.company_id, sv.status, salary, 
		vc.date_last_change, vc.rej_text
FROM vacancies JOIN companies on (vacancies.company_id=companies.id)
	JOIN vacancies_to_candidates vc on (vacancies.id = vc.vacancy_id)
	JOIN platforms ON (vacancies.platform_id = platforms.id)
	JOIN status_for_vacs sv on vc.status = sv.id
	WHERE vc.candidate_id=$1`, ref.ViewCandidate.Id)
	if err != nil {
		logs.ErrorLog(err, "")
	}

	if len(ref.ViewCandidate.Vacancies) > 0 {
		ref.Status.CompId, _ = ref.ViewCandidate.Vacancies[0]["company_id"].(int32)
		ref.Status.CompName, _ = ref.ViewCandidate.Vacancies[0]["name"].(string)
		ref.Status.Comments, _ = ref.ViewCandidate.Vacancies[0]["rej_text"].(string)
		ref.Status.VacStat, _ = ref.ViewCandidate.Vacancies[0]["status"].(string)
		ref.Status.Date, _ = ref.ViewCandidate.Vacancies[0]["date_last_change"].(time.Time)
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

type SelectedUnit struct {
	Id    int32  `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type SelectedUnits []*SelectedUnit

func (s *SelectedUnits) GetFields(columns []dbEngine.Column) []interface{} {
	p := &SelectedUnit{}
	r := make([]interface{}, 0)
	for _, col := range columns {
		switch col.Name() {
		case "id":
			r = append(r, &p.Id)
		case "label":
			r = append(r, &p.Label)
		case "value":
			r = append(r, &p.Value)
		default:
			logs.ErrorLog(dbEngine.ErrNotFoundColumn{Column: col.Name()}, "SelectedUnits. GetFields")
		}
	}

	*s = append(*s, p)

	return r
}

type ResList struct {
	Count, Page, TotalPage int
	CurrentPage            int           `json:"currentPage"`
	PerPage                int           `json:"perPage"`
	Platforms              SelectedUnits `json:"platforms"`
	Seniority              SelectedUnits `json:"seniority"`
}

func NewResList(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, pageNum int) *ResList {
	return &ResList{
		Page:        pageNum,
		CurrentPage: pageNum,
		Count:       10 * pageItem,
		TotalPage:   10,
		PerPage:     pageItem,
		Platforms:   getPlatforms(ctx, DB),
		Seniority:   getSeniorities(ctx, DB),
	}
}

func getVacToCand(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	err := DB.Conn.SelectAndScanEach(ctx,
		nil,
		&res,
		`select v.status as id, s.status  as label, lower(s.status) as value
        from vacancies_to_candidates v join status_for_vacs s on s.id = status)
`,
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getStatusVac(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	statUses, _ := db.NewStatus_for_vacs(DB)

	err := statUses.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "status as label", "LOWER(status) as value"),
		dbEngine.OrderBy("order_num"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getStatuses(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	statUses, _ := db.NewStatuses(DB)

	err := statUses.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "status as label", "LOWER(status) as value"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getTags(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	statUses, _ := db.NewTags(DB)

	err := statUses.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getRejectReason(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	statUses, _ := db.NewTags(DB)

	err := statUses.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
		dbEngine.WhereForSelect("parent_id"),
		dbEngine.ArgsForSelect(3),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getLocations(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	statUses, _ := db.NewLocation_for_vacancies(DB)

	err := statUses.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getSeniorities(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	statUses, _ := db.NewSeniorities(DB)

	err := statUses.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "nazva as label", "LOWER(nazva) as value"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getPlatforms(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	platforms, _ := db.NewPlatforms(DB)

	err := platforms.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "nazva as label", "LOWER(nazva) as value"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getRecruters(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	users, _ := db.NewUsers(DB)

	err := users.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getCompanies(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, opt ...dbEngine.BuildSqlOptions) (res SelectedUnits) {
	company, _ := db.NewCompanies(DB)

	opt = append(opt, dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"))
	err := company.SelectAndScanEach(ctx,
		nil,
		&res,
		opt...,
	)
	if err != nil {
		logs.ErrorLog(err, "	SelectSelfScanEach")
	}

	return
}
