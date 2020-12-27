// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

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

type ResCandidates struct {
	Count, Page int
	CurrentPage int                    `json:"currentPage"`
	PerPage     int                    `json:"perPage"`
	Candidates  []*db.CandidatesFields `json:"candidates"`
	Company     []*db.CompaniesFields  `json:"company"`
	Platforms   []*db.PlatformsFields  `json:"platforms"`
	Recruiter   []string               `json:"recruiter"`
	Statuses    []*db.StatusesFields   `json:"statuses"`
}

type selectOpt struct {
	Companies     []*db.CompaniesFields `json:"companies"`
	Platforms     []*db.PlatformsFields `json:"platforms"`
	Recruiters    []string              `json:"recruiters"`
	Statuses      []*db.StatusesFields  `json:"candidateStatus"`
	Location      []*db.StatusesFields  `json:"location"`
	RejectReasons []*db.StatusesFields  `json:"reject_reasons"`
	RejectTag     []*db.StatusesFields  `json:"reject_tag"`
	Seniorities   []*db.StatusesFields  `json:"seniorities"`
	Tags          []*db.StatusesFields  `json:"tags"`
	VacancyStatus []*db.StatusesFields  `json:"vacancyStatus"`
}
type ViewCandidate struct {
	Candidates []*db.CandidatesFields `json:"0"`
	SelectOpt  selectOpt              `json:"select"`
	Statuses   []*db.StatusesFields   `json:"statusesCandidate"`
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
		Candidates:  make([]*db.CandidatesFields, pageItem),
		Recruiter:   []string{},
		// Company:    make([]*db.CompaniesFields, pageItem),
	}
	i := 0
	err := table.SelectSelfScanEach(ctx,
		func(record *db.CandidatesFields) error {
			res.Candidates[i] = record
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

	company, _ := db.NewCompanies(DB)

	err = company.SelectSelfScanEach(ctx,
		func(record *db.CompaniesFields) error {
			res.Company = append(res.Company, record)

			return nil
		},
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	platforms, _ := db.NewPlatforms(DB)

	err = platforms.SelectSelfScanEach(ctx,
		func(record *db.PlatformsFields) error {
			res.Platforms = append(res.Platforms, record)

			return nil
		},
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	statUses, _ := db.NewStatuses(DB)

	err = statUses.SelectSelfScanEach(ctx,
		func(record *db.StatusesFields) error {
			res.Statuses = append(res.Statuses, record)

			return nil
		},
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
	res := ViewCandidate{
		// Page:        1,
		// CurrentPage: 1,
		// PerPage:     pageItem,
		Candidates: make([]*db.CandidatesFields, 1),
		// Recruiter:   []string{},
		// Company:    make([]*db.CompaniesFields, pageItem),
	}

	res.Candidates[0] = table.Record
	company, _ := db.NewCompanies(DB)

	err = company.SelectSelfScanEach(ctx,
		func(record *db.CompaniesFields) error {
			res.SelectOpt.Companies = append(res.SelectOpt.Companies, record)

			return nil
		},
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	platforms, _ := db.NewPlatforms(DB)

	err = platforms.SelectSelfScanEach(ctx,
		func(record *db.PlatformsFields) error {
			res.SelectOpt.Platforms = append(res.SelectOpt.Platforms, record)

			return nil
		},
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	statUses, _ := db.NewStatuses(DB)

	err = statUses.SelectSelfScanEach(ctx,
		func(record *db.StatusesFields) error {
			res.Statuses = append(res.Statuses, record)

			return nil
		},
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	return res, nil
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
		dbEngine.ArgsForSelect(u.Name,
			u.Salary,
			u.Email,
			u.Mobile,
			u.Skype,
			u.Link,
			u.Linkedin,
			u.Str_companies,
			u.Status,
			u.Tag_id,
			u.Comments,
			u.Date,
			u.Recruter_id,
			u.Text_rezume,
			u.Sfera,
			u.Experience,
			u.Education,
			u.Language,
			u.Zapoln_profile,
			u.File,
			u.Avatar,
			u.Seniority_id,
			u.Date_follow_up,
		),
	)
	if err != nil {
		return nil, err
	}

	return i, nil
}
