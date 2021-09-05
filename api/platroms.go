// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"sort"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/uppeople/db"
	"github.com/valyala/fasthttp"
)

type ResPlatforms struct {
	*ResList
	Platforms []*db.PlatformsFields `json:"platforms"`
}

func HandleAddPlatform(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	p, ok := ctx.UserValue(apis.JSONParams).(*db.PlatformsFields)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	platforms, _ := db.NewPlatforms(DB)

	i, err := platforms.Insert(ctx,
		dbEngine.ColumnsForSelect("name"),
		dbEngine.ArgsForSelect(p.Name),
	)
	if err != nil {
		return createErrResult(err)
	}

	return createResult(i)
}

func HandleAllPlatforms(ctx *fasthttp.RequestCtx) (interface{}, error) {
	id, ok := ctx.UserValue(ParamPageNum.Name).(int)
	if !ok {
		id = 1
	}

	res := ResPlatforms{
		ResList: NewResList(id),
	}

	platformsMap := db.GetPlatformsAsMap()
	platformNames := make([]string, len(platformsMap))
	i := 0
	for k := range platformsMap {
		platformNames[i] = k
		i++
	}
	// TODO: write custom sorting algorithm for string slice with upper case and lower case first characters
	sort.Strings(platformNames)
	for num, platform := range platformNames {
		if (id-1)*pageItem <= num && num < id*pageItem {
			if platformVal, ok := platformsMap[platform]; ok {
				res.Platforms = append(res.Platforms, &platformVal)
			}
		}
	}

	if len(platformNames) < pageItem {
		res.ResList.TotalPage = 1
	} else {
		res.ResList.TotalPage = len(platformNames) / pageItem

	}
	res.ResList.Count = len(platformNames)

	return res, nil

}
