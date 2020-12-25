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

type ResCandidates struct {
	Count, Page, CurrentPage, PerPage int
	Candidates                        []*db.CandidatesFields
	Company                           []*db.CompaniesFields
}

const pageItem = 14

func HandleAllCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewCandidates(DB)
	res := ResCandidates{
		Candidates: make([]*db.CandidatesFields, pageItem),
		Company:    make([]*db.CompaniesFields, pageItem),
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

	i = 0
	err = company.SelectSelfScanEach(ctx,
		func(record *db.CompaniesFields) error {
			res.Company[i] = record
			i++
			if i == pageItem {
				return errLimit
			}

			return nil
		},
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	return res, nil
}

func HandleAddCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	user, err := prepareRequest(ctx)
	if err != nil {
		return nil, err
	}

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
