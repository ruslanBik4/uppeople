// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/ruslanBik4/uppeople/db"
	"github.com/valyala/fasthttp"
)

type ResPlatforms struct {
	*ResList
	Platforms []*db.PlatformsFields `json:"platforms"`
}

func HandleAllPlatforms(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	offset := 0
	id, ok := ctx.UserValue(ParamPageNum.Name).(int)
	if ok && id > 1 {
		offset = id * pageItem
	}

	platforms, _ := db.NewPlatforms(DB)
	res := ResPlatforms{
		ResList: NewResList(id),
	}

	optionsCount := []dbEngine.BuildSqlOptions{
		dbEngine.ColumnsForSelect("count(*)"),
	}
	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy("nazva"),
		dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	}
	err := platforms.SelectSelfScanEach(ctx,
		func(record *db.PlatformsFields) error {
			res.Platforms = append(res.Platforms, record)
			return nil
		},
		options...,
	)
	if err != nil {
		return createErrResult(err)
	}

	if len(res.Platforms) < pageItem {
		res.ResList.TotalPage = 1
		res.ResList.Count = len(res.Platforms)
	} else {
		err = platforms.SelectOneAndScan(ctx,
			&res.ResList.Count,
			optionsCount...)
		if err != nil {
			logs.ErrorLog(err, "count")
		} else {
			res.ResList.TotalPage = res.ResList.Count / pageItem
		}
	}

	return res, nil

}
