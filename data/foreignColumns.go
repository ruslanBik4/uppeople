// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/views/templates/forms"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

func GetForeignOptions(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, colDec *forms.ColumnDecor, id interface{}) {
	if colDec.Name() == "id_polymers" {
		if ctx.UserValue("html") == nil {
			colDec.Suggestions = "/search/list"
		} else {
			colDec.Suggestions = "api/v1/search/list"

		}
		return
		// 	todo: mv this into DB
	}

	if strings.HasPrefix(colDec.Name(), "id_") {
		table, ok := DB.Tables[strings.TrimPrefix(colDec.Name(), "id_")]
		if ok {
			colDec.SelectOptions = make(map[string]string)
			// todo will decoded about clear value of field
			// if id != nil && !colDec.Required() {
			// 	colDec.SelectOptions["disabled"] = ""
			// }

			name := GetNameOfTitleColumn(table, ctx.UserValue(ParamsLang.Name))

			err := table.SelectAndRunEach(ctx,
				func(values []interface{}, columns []dbEngine.Column) error {
					colDec.SelectOptions[values[1].(string)] = strconv.Itoa(int(values[0].(int32)))
					return nil
				},
				dbEngine.ColumnsForSelect("id", name),
				dbEngine.OrderBy(name),
			)
			if err != nil {
				logs.ErrorLog(err, "")
			}

			if !colDec.IsSlice && id != nil {
				colDec.Value = fmt.Sprintf("%d", id)
			}
		}
	}
	getSystemField(ctx, DB, colDec, id)
}

func getSystemField(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, colDec *forms.ColumnDecor, id interface{}) {
	switch colDec.Name() {
	case "tablename":
		colDec.SelectOptions = make(map[string]string, len(DB.Tables))
		for name, table := range DB.Tables {
			label := table.Comment()
			if strings.TrimSpace(label) == "" {
				label = name
			}
			colDec.SelectOptions[label] = name
		}
	case "columns":
		colDec.SelectOptions = make(map[string]string)
		for _, col := range DB.Tables["items"].Columns() {
			label := col.Comment()
			if strings.TrimSpace(label) == "" {
				label = col.Name()
			} else if p := strings.Split(label, "{"); len(p) > 1 {
				label = p[0]
			}
			colDec.SelectOptions[label] = col.Name()
		}
	default:

	}
}

func GetForeignName(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, col dbEngine.Column, val interface{}) interface{} {
	if val != nil && strings.HasPrefix(col.Name(), "id_") {
		table, ok := DB.Tables[strings.TrimPrefix(col.Name(), "id_")]
		if ok {

			name := GetNameOfTitleColumn(table, ctx.UserValue(ParamsLang.Name))
			if strings.HasPrefix(col.Type(), "_") {
				res := make([]string, 0)
				err := DB.Conn.SelectOneAndScan(ctx,
					&res,
					fmt.Sprintf("select array_agg(%s) from %s where id =ANY($1)", name, table.Name()),
					val,
				)
				if err != nil {
					logs.ErrorLog(err, "%s=%v", name, val)
					return nil
				}
				return res
			}

			res := ""
			err := table.SelectOneAndScan(ctx,
				&res,
				dbEngine.ColumnsForSelect(name),
				dbEngine.WhereForSelect("id"),
				dbEngine.ArgsForSelect(val),
			)
			if err != nil {
				logs.ErrorLog(err, "%s=%v", name, val)
				return nil
			}
			return res
		}
	}

	return nil
}

var names = []string{
	"name",
	"title",
	"desc",
	"description",
}

func GetNameOfTitleColumn(table dbEngine.Table, lang interface{}) string {
	for _, name := range names {
		col := table.FindColumn(name)
		if col != nil {
			return GetNameAccordingLang(table, name, lang)
		}
	}
	for _, col := range table.Columns() {
		if col.Name() != "id" {
			return GetNameAccordingLang(table, col.Name(), lang)
		}
	}

	return ""
}

func GetNameAccordingLang(table dbEngine.Table, name string, lang interface{}) string {

	if l, ok := lang.(string); ok && (table.FindColumn(name+"_"+l) != nil) {
		return name + "_" + l
	}

	return name
}
