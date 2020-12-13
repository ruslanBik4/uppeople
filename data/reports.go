// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"go/types"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"
)

type reportHeader struct {
	Name  string
	Title string `json:"title"`
	Typ   string `json:"type"`
	Total string
}

type ReportJSON struct {
	Header   []reportHeader
	Data     [][]interface{}
	Settings map[string]interface{}
	table    dbEngine.Table
}

func NewReportJSON(table dbEngine.Table) *ReportJSON {
	r := &ReportJSON{
		Header:   make([]reportHeader, len(table.Columns())),
		Settings: make(map[string]interface{}),
		table:    table,
	}
	for i, col := range table.Columns() {
		comment := col.Comment()
		if comment == "" {
			comment = col.Name()
		}
		typ := col.Type()
		switch {
		case (col.BasicTypeInfo() & types.IsInteger) != 0:
			typ = "integer"
		case (col.BasicTypeInfo() & types.IsFloat) != 0:
			typ = "float"
		case col.BasicTypeInfo() == types.IsString:
			typ = "string"
		case col.BasicTypeInfo() == types.IsBoolean:
			typ = "bool"
		}

		r.Header[i] = reportHeader{
			Name:  col.Name(),
			Title: comment,
			Typ:   typ,
			Total: "",
		}
	}

	return r
}

func (r *ReportJSON) getRoute() *apis.ApiRoute {
	rReport := &apis.ApiRoute{
		Desc: "struct table '" + r.table.Name() + "' data",
		// Params:    make([]apis.InParam, 0),
	}

	rReport.Fnc = func(ctx *fasthttp.RequestCtx) (interface{}, error) {

		err := r.table.SelectAndRunEach(ctx,
			func(values []interface{}, columns []dbEngine.Column) error {
				r.Data = append(r.Data, values)
				return nil
			})
		if err != nil {
			return nil, err
		}

		return r, nil
	}

	return rReport
}
