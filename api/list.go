// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
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
			CompId:       0,
			Recruiter:    "",
			DateFollowUp: record.Date_follow_up,
		},
	}
	if record.Recruter_id.Valid {
		recTable, _ := db.NewUsers(DB)

		err := recTable.SelectOneAndScan(ctx,
			&ref.Recruiter,
			dbEngine.ColumnsForSelect("name"),
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(record.Recruter_id.Int64),
		)
		if err != nil {
			logs.ErrorLog(err, "recTable.SelectOneAndScan")
		}
		ref.Status.Recruiter = ref.Recruiter
	}

	tagTable, _ := db.NewTags(DB)

	err := tagTable.SelectOneAndScan(ctx,
		ref.Tags,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(record.Tag_id),
	)
	if err != nil {
		logs.ErrorLog(err, "recTable.SelectOneAndScan")
	} else {
		ref.TagName = ref.Tags.Name
		ref.TagColor = ref.Tags.Color
	}

	for _, s := range seniors {
		if s.Id == int32(ref.Seniority_id.Int64) {
			ref.Seniority = s.Label
		}
	}

	for _, p := range platforms {
		if p.Id == int32(record.Platform_id.Int64) {
			ref.Platform = p.Label
			ref.ViewCandidate.Platform = p
		}
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

func NewResList(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) *ResList {
	return &ResList{
		Page:        1,
		CurrentPage: 1,
		PerPage:     pageItem,
		Platforms:   getPlatforms(ctx, DB),
		Seniority:   getSeniorities(ctx, DB),
	}
}

type selectOpt struct {
	Companies     SelectedUnits `json:"companies"`
	Platforms     SelectedUnits `json:"platforms"`
	Recruiters    SelectedUnits `json:"recruiters"`
	Statuses      SelectedUnits `json:"candidateStatus"`
	Location      SelectedUnits `json:"location"`
	RejectReasons SelectedUnits `json:"reject_reasons"`
	RejectTag     SelectedUnits `json:"reject_tag"`
	Recruiter     SelectedUnits `json:"recruiter"`
	Seniorities   SelectedUnits `json:"seniorities"`
	Tags          SelectedUnits `json:"tags"`
	VacancyStatus SelectedUnits `json:"vacancyStatus"`
	// vacancies
}

func getVacToCand(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	vacCand, _ := db.NewVacancies_to_candidates(DB)
	err := vacCand.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "status as label", "LOWER(status) as value"),
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

func getCompanies(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) (res SelectedUnits) {
	company, _ := db.NewCompanies(DB)

	err := company.SelectAndScanEach(ctx,
		nil,
		&res,
		dbEngine.ColumnsForSelect("id", "name as label", "LOWER(name) as value"),
	)
	if err != nil {
		logs.ErrorLog(err, "	SelectSelfScanEach")
	}

	return
}
