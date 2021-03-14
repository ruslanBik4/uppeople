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
	if err != nil && err != pgx.ErrNoRows {
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
		`select v.id, 
		concat(companies.name, ' ("', platforms.nazva, '")') as name, 
		concat(companies.name, ' ("', platforms.nazva, '")') as label, 
		LOWER(CONCAT(companies.name, ' ("', platforms.nazva , '")')) as value, 
		user_ids, 
		platform_id,
		CONCAT(platforms.nazva, ' ("', (select nazva from seniorities where id=seniority_id), '")') as platform,
		companies, sv.id as status_id, v.company_id, sv.status, salary, 
		coalesce(vc.date_last_change, vc.date_create, $3) as date_last_change, vc.rej_text, sv.color
FROM vacancies v JOIN companies on (v.company_id=companies.id)
	LEFT JOIN vacancies_to_candidates vc on (v.id = vc.vacancy_id AND vc.candidate_id=$1)
	JOIN platforms ON (v.platform_id = platforms.id)
	JOIN status_for_vacs sv on coalesce(vc.status, 1) = sv.id
	WHERE (vc.candidate_id=$1 OR v.id = ANY($2)) AND v.status!=1
    order by date_last_change desc
`,
		ref.ViewCandidate.Id,
		ref.ViewCandidate.CandidatesFields.Vacancies,
		"2020-01-01",
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
        from vacancies_to_candidates v join status_for_vacs s on (s.id = v.status)
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
		dbEngine.WhereForSelect("parent_id"),
		dbEngine.ArgsForSelect(0),
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
		dbEngine.OrderBy("name"),
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
		// dbEngine.OrderBy("nazva"),
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
		dbEngine.OrderBy("nazva"),
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
		dbEngine.OrderBy("name"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getCompanies(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, opt ...dbEngine.BuildSqlOptions) (res SelectedUnits) {
	company, _ := db.NewCompanies(DB)

	err := company.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
		dbEngine.OrderBy("name"),
	)
	if err != nil {
		logs.ErrorLog(err, "	SelectSelfScanEach")
	}

	return
}
