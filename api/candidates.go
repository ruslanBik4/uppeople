// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"

	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type CandidateDTO struct {
	*db.CandidatesFields
}

func (c *CandidateDTO) GetValue() interface{} {
	panic("implement me")
}

func (c *CandidateDTO) NewValue() interface{} {
	return &db.CandidatesFields{}
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

func HandleAddCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*db.CandidatesFields)
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
	user, err := prepareRequest(ctx)
	if err != nil {
		return nil, err
	}

	return i, nil

	var v map[string]interface{}
	err = jsoniter.Unmarshal(ctx.Request.Body(), &v)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	req := &fasthttp.Request{}
	ctx.Request.CopyTo(req)

	req.URI().SetPath("api/main/allCandidates/1")
	nameCan := v["name"].(string)
	req.SetBodyString(`{"search":"` + nameCan + `","dateFrom":"","dateTo":"","selectCompanies":[],"selectPlatforms":[],"selectStatuses":[],"selectRecruiter":"","selectTag":[],"selectReason":[],"mySent":false,"dateFromAllCandidates":"","dateToAllCandidates":"","dateFromSentCandidates":"","dateToSentCandidates":"","dateFromFreelancersCandidates":"","dateToFreelancersCandidates":"","dateFollowUpFrom":"","dateFollowUpTo":"","selectSeniority":[]}`)
	resp := &fasthttp.Response{}
	err = doRequest(req, resp, user.Host)
	if err != nil {
		return nil, err
	}

	if bytes.Contains(resp.Header.ContentType(), []byte("json")) {
		err = jsoniter.Unmarshal(resp.Body(), &v)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal")
		}

		c, ok := v["Count"].(float64)
		if !ok {
			logs.DebugLog("%T", v["Count"])
		} else if c > 0 {
			return nameCan + " already exists", apis.ErrWrongParamsList
		} else {
			logs.DebugLog("%f", v["Count"])
			return nil, doRequest(&ctx.Request, &ctx.Response, user.Host)
		}
	}

	return nil, nil
}
