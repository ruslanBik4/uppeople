// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/ruslanBik4/uppeople/db"
	"github.com/valyala/fasthttp"
)

func HandleAddAvatar(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	input, ok := ctx.UserValue("file").([]*multipart.FileHeader)
	if !ok {
		return map[string]string{
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	_, bytea, err := readByteA(input)
	table, _ := db.NewCandidates(DB)
	args := []interface{}{
		string(bytea[0]),
		id,
	}
	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect("avatar"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	if i > 0 {
		toLogCandidate(ctx, DB, id, "avatar", CODE_LOG_UPDATE)
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}

func HandleAddLogo(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	input, ok := ctx.UserValue("file").([]*multipart.FileHeader)
	if !ok {
		return map[string]string{
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	names, bytea, err := readByteA(input)
	if err != nil {
		return createErrResult(err)
	}
	logo, err := HandleSaveImg(ctx, bytea[0], names[0])
	if err == nil {
		table, _ := db.NewCompanies(DB)
		args := []interface{}{
			logo,
			id,
		}
		_, err = table.Update(ctx,
			dbEngine.ColumnsForSelect("logo"),
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(args...),
		)

	}
	if err != nil {
		return createErrResult(err)
	}

	toLogCompany(ctx, DB, id, "logo", CODE_LOG_UPDATE)
	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}

func readByteA(fHeaders []*multipart.FileHeader) ([]string, [][]byte, error) {
	bytea := make([][]byte, len(fHeaders))
	names := make([]string, len(fHeaders))
	for i, fHeader := range fHeaders {

		f, err := fHeader.Open()
		if err != nil {
			logs.DebugLog(err, fHeader)
			return nil, nil, errors.Wrap(err, fHeader.Filename)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			logs.DebugLog(err, fHeader)
			return nil, nil, errors.Wrap(err, "read "+fHeader.Filename)
		}

		bytea[i] = b
		names[i] = fHeader.Filename
	}

	return names, bytea, nil
}
