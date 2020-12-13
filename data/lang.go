// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"runtime"
	"strings"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

func Translate(ctx *fasthttp.RequestCtx, name string) string {
	lang, ok := ctx.UserValue("lang").(string)
	if !ok || lang == "eng" || name == "id" {
		return name
	}

	name = strings.TrimSpace(name)
	if name == "" {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			logs.DebugLog("name empty ", runtime.FuncForPC(pc).Name(), file, line)
		}
		return ""
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return name
	}

	const sqlFindTranslate = `select translation from dictionary d where d.name=$1
     and id_languages = (select id from languages l where l.name=$2)`

	var trans string
	err := DB.Conn.SelectOneAndScan(ctx, &trans, sqlFindTranslate, name, lang)
	if err != nil || trans == "" {
		logs.ErrorLog(err, "no translation", name)
		return name
	}

	return trans
}
