// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

type SelectedUnit struct {
	Id    int32  `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type CandidateDTO struct {
	*db.CandidatesFields
	Comment           string         `json:"comment"`
	Date              string         `json:"date"`
	Phone             string         `json:"phone"`
	Resume            string         `json:"resume"`
	SelectPlatform    SelectedUnit   `json:"selectPlatform"`
	SelectSeniority   SelectedUnit   `json:"selectSeniority"`
	SelectedTag       SelectedUnit   `json:"selectedTag"`
	SelectedVacancies []SelectedUnit `json:"selectedVacancies"`
}

func (c *CandidateDTO) GetValue() interface{} {
	return c
}

func (c *CandidateDTO) NewValue() interface{} {
	return &CandidateDTO{CandidatesFields: &db.CandidatesFields{}}
}

type CandidateView struct {
	*db.CandidatesFields
	Platform  string              `json:"platform"`
	Platforms *db.PlatformsFields `json:"platforms,omitempty"`
	Seniority string              `json:"seniority"`
	TagName   string              `json:"tag_name"`
	Tags      *db.TagsFields      `json:"tags,omitempty"`
	TagColor  string              `json:"tag_color"`
	Recruiter string              `json:"recruiter"`
	// status
}
type ResCandidates struct {
	Count, Page int
	CurrentPage int                   `json:"currentPage"`
	PerPage     int                   `json:"perPage"`
	Candidates  []CandidateView       `json:"candidates"`
	Company     []*db.CompaniesFields `json:"company"`
	Platforms   []*db.PlatformsFields `json:"platforms"`
	Recruiter   []string              `json:"recruiter"`
	Statuses    []*db.StatusesFields  `json:"statuses"`
}

type selectOpt struct {
	Companies     []*db.CompaniesFields              `json:"companies"`
	Platforms     []*db.PlatformsFields              `json:"platforms"`
	Recruiters    []string                           `json:"recruiters"`
	Statuses      []*db.StatusesFields               `json:"candidateStatus"`
	Location      []*db.Location_for_vacanciesFields `json:"location"`
	RejectReasons []string                           `json:"reject_reasons"`
	RejectTag     []*db.TagsFields                   `json:"reject_tag"`
	Seniorities   []*db.SenioritiesFields            `json:"seniorities"`
	Tags          []*db.TagsFields                   `json:"tags"`
	VacancyStatus []*db.Status_for_vacsFields        `json:"vacancyStatus"`
}
type VacanciesDTO struct {
	*db.VacanciesFields
	Platforms *db.PlatformsFields `json:"platforms"`
}
type StatusesCandidate struct {
	Candidate_id     int32                     `json:"candidate_id"`
	Company          *db.CompaniesFields       `json:"company"`
	Company_id       int32                     `json:"company_id"`
	Date_create      time.Time                 `json:"date_create"`
	Date_last_change time.Time                 `json:"date_last_change"`
	Id               int32                     `json:"id"`
	Notice           string                    `json:"notice"`
	Rating           string                    `json:"rating"`
	Rej_text         string                    `json:"rej_text"`
	Status           int32                     `json:"status"`
	Status_vac       *db.Status_for_vacsFields `json:"status_vac"`
	User_id          int32                     `json:"user_id"`
	Vacancy          VacanciesDTO              `json:"vacancy"`
	Vacancy_id       int32                     `json:"vacancy_id"`
}
type ViewCandidates struct {
	Candidates []CandidateView     `json:"0"`
	SelectOpt  selectOpt           `json:"select"`
	Statuses   []StatusesCandidate `json:"statusesCandidate"`
}

const pageItem = 15

func HandleAllCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewCandidates(DB)
	res := ResCandidates{
		Page:        1,
		CurrentPage: 1,
		PerPage:     pageItem,
		Candidates:  make([]CandidateView, pageItem),
		Recruiter:   []string{},
		Company:     getCompanies(ctx, DB),
		Platforms:   getPlatforms(ctx, DB),
		Statuses:    getStatUses(ctx, DB),
	}

	tags := getTags(ctx, DB)
	seniors := getSeniorities(ctx, DB)
	recTable, _ := db.NewUsers(DB)
	i := 0
	err := table.SelectSelfScanEach(ctx,
		func(record *db.CandidatesFields) error {

			cV := CandidateView{
				CandidatesFields: record,
			}
			if record.Recruter_id.Valid {
				err := recTable.SelectOneAndScan(ctx,
					&cV.Recruiter,
					dbEngine.ColumnsForSelect("name"),
					dbEngine.WhereForSelect("id"),
					dbEngine.ArgsForSelect(record.Recruter_id.Int64),
				)
				if err != nil {
					logs.ErrorLog(err, "recTable.SelectOneAndScan")
				}
			}
			for _, tag := range tags {
				if tag.Id == record.Tag_id {
					cV.TagName = tag.Name
					cV.TagColor = tag.Color
				}
			}
			for _, s := range seniors {
				if s.Id == cV.Seniority_id.Int64 {
					cV.Seniority = s.Nazva.String
				}
			}
			for _, p := range res.Platforms {
				if p.Id == record.Platform_id.Int64 {
					cV.Platform = p.Nazva.String
				}
			}

			res.Candidates[i] = cV
			i++
			if i == pageItem {
				return errLimit
			}

			return nil
		},
		dbEngine.OrderBy("date desc"),
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	return res, nil
}

func HandleViewCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewCandidates(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(ctx.UserValue("id")),
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}
	res := ViewCandidates{
		Candidates: []CandidateView{
			{
				CandidatesFields: table.Record,
			},
		},
		SelectOpt: selectOpt{
			Companies:     getCompanies(ctx, DB),
			Platforms:     getPlatforms(ctx, DB),
			Recruiters:    make([]string, 0),
			Statuses:      getStatUses(ctx, DB),
			Location:      getLocations(ctx, DB),
			RejectReasons: []string{},
			RejectTag:     getTags(ctx, DB),
			Seniorities:   getSeniorities(ctx, DB),
			Tags:          getTags(ctx, DB),
			VacancyStatus: getStatusVac(ctx, DB),
		},
		Statuses: []StatusesCandidate{
			{
				Candidate_id: int32(table.Record.Id),
				Company:      &db.CompaniesFields{},
				Status_vac:   &db.Status_for_vacsFields{},
				Vacancy: VacanciesDTO{
					&db.VacanciesFields{},
					&db.PlatformsFields{},
				},
			},
		},
	}

	for _, p := range res.SelectOpt.Platforms {
		if p.Id == table.Record.Platform_id.Int64 {
			res.Candidates[0].Platforms = p
		}
	}

	for _, tag := range res.SelectOpt.Tags {
		if tag.Id == table.Record.Tag_id {
			res.Candidates[0].Tags = tag
		}
	}

	return res, nil
}

func getStatUses(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.StatusesFields {
	statUses, _ := db.NewStatuses(DB)
	res := make([]*db.StatusesFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.StatusesFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getStatusVac(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.Status_for_vacsFields {
	statUses, _ := db.NewStatus_for_vacs(DB)
	res := make([]*db.Status_for_vacsFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.Status_for_vacsFields) error {
			res = append(res, record)

			return nil
		},
		dbEngine.OrderBy("order_num"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getTags(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.TagsFields {
	statUses, _ := db.NewTags(DB)
	res := make([]*db.TagsFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.TagsFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getLocations(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.Location_for_vacanciesFields {
	statUses, _ := db.NewLocation_for_vacancies(DB)
	res := make([]*db.Location_for_vacanciesFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.Location_for_vacanciesFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getSeniorities(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.SenioritiesFields {
	statUses, _ := db.NewSeniorities(DB)
	res := make([]*db.SenioritiesFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.SenioritiesFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getPlatforms(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.PlatformsFields {
	platforms, _ := db.NewPlatforms(DB)
	res := make([]*db.PlatformsFields, 0)

	err := platforms.SelectSelfScanEach(ctx,
		func(record *db.PlatformsFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getCompanies(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.CompaniesFields {
	company, _ := db.NewCompanies(DB)
	companies := make([]*db.CompaniesFields, 0)

	err := company.SelectSelfScanEach(ctx,
		func(record *db.CompaniesFields) error {
			companies = append(companies, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	SelectSelfScanEach")
	}

	return companies
}

func HandleAddCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*CandidateDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	logs.DebugLog(u)
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewCandidates(DB)
	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(
			"name",
			"platform_id",
			"salary",
			"email",
			"mobile",
			"skype",
			"link",
			"linkedin",
			"str_companies",
			"status",
			"tag_id",
			"comments",
			"date",
			"recruter_id",
			"text_rezume",
			"sfera",
			"experience",
			"education",
			"language",
			"zapoln_profile",
			"file",
			"avatar",
			"seniority_id",
			"date_follow_up",
		),
		dbEngine.ArgsForSelect(
			u.Name,
			u.SelectPlatform.Id,
			u.Salary,
			u.Email,
			u.Phone,
			u.Skype,
			u.Link,
			u.Linkedin,
			u.Str_companies,
			u.Status,
			u.SelectedTag.Id,
			u.Comment,
			u.Date,
			auth.GetUserData(ctx).Id,
			u.Resume,
			u.Sfera,
			u.Experience,
			u.Education,
			u.Language,
			u.Zapoln_profile,
			u.File,
			u.Avatar,
			u.SelectSeniority.Id,
			u.Date_follow_up,
		),
	)
	if err != nil {
		return nil, err
	}

	return i, nil
}
