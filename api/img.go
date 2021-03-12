// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
)

func HandleGetImg(ctx *fasthttp.RequestCtx) (interface{}, error) {
	path, ok := ctx.UserValue(apis.ChildRoutePath).(string)
	if !ok {
		return "", apis.ErrWrongParamsList
	}
	b, err := ioutil.ReadFile(filepath.Join("img", path))
	if err != nil {
		logs.ErrorLog(err, "read file")
		b, err = ioutil.ReadFile(filepath.Join("img", "companies_logo/no_logo.png"))
		if err != nil {
			return createErrResult(err)
		}
	}

	download(ctx, b, path)

	return nil, nil
}

func HandleImg(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	photos, ok := DB.Tables["photos"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "photos"}
	}

	b := make([]byte, 0)
	name := ""
	row := []interface{}{
		&b, &name,
	}
	err := photos.SelectOneAndScan(ctx,
		row,
		dbEngine.ColumnsForSelect("blob", "name"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	download(ctx, b, name)

	return nil, nil
}

// todo change to ioReader
func HandleSaveImg(ctx *fasthttp.RequestCtx, data []byte, name string) (string, error) {
	path, ok := ctx.UserValue(apis.ChildRoutePath).(string)
	if !ok {
		return "", apis.ErrWrongParamsList
	}

	fullPath := filepath.Join("img", "companies_logo", path)
	err := os.Mkdir(fullPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return "", err
	}

	fileName := filepath.Join(fullPath, name)
	err = ioutil.WriteFile(fileName, data, os.ModePerm)
	if err != nil {
		return "", err
	}

	return filepath.Join("companies_logo", path, name), nil
}

func download(ctx *fasthttp.RequestCtx, body []byte, filename string) {
	if body != nil {
		ctx.Response.SetBody(body)
	}
	// Send the headers
	ctx.Response.Header.Set("Content-Description", "File Transfer")
	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(body)))
	ctx.Response.Header.Set("Last-Modified", time.Now().String())
	if bytes.HasPrefix(body, []byte("<?xml")) ||
		bytes.HasPrefix(body, []byte("<svg")) {
		ctx.Response.Header.SetContentType("image/svg+xml")
		filename += ".svg"
	} else {
		ctx.Response.Header.SetContentType("application/octet-stream")
		ctx.Response.Header.Set("Content-Transfer-Encoding", "binary")
		cType := http.DetectContentType(body)
		ctx.SetContentType(cType)
		if strings.HasSuffix(cType, "pdf") {
			filename += ".pdf"
		}
	}

	ctx.Response.Header.Set("Content-Disposition", "attachment; filename="+filename)
	if ctx.UserValue("v") != nil {
		ctx.Response.Header.Set("Cache-Control", "max-age=0, no-cache, no-store")
		ctx.Response.Header.Set("Pragma", "no-cache")
		ctx.Response.Header.Set("Expires", "Wed, 11 Jan 1984 05:00:00 GMT")
	} else {
		ctx.Response.Header.Set("Cache-Control", "must-revalidate")
	}
}
