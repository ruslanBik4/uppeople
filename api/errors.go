// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"regexp"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
)

var (
	regKeyWrong   = regexp.MustCompile(`Key\s+\((\w+)\)=\((.+)\)([^.]+)`)
	regDublicated = regexp.MustCompile(`duplicate key value violates unique constraint "(\w*)"`)

	TagIdDateNotUpdatedOnSendCV = errors.New("TagId and date not updated on SendCV")
)

// todo: add code NoContent
func createErrResult(err error) (interface{}, error) {
	if err == pgx.ErrNoRows {
		return nil, err
	}
	msg := err.Error()
	e, ok := errors.Cause(err).(*pgconn.PgError)
	if ok {
		msg = e.Detail
	}

	if s := regKeyWrong.FindStringSubmatch(msg); len(s) > 0 {
		return map[string]string{
			s[1]: "`" + s[2] + "`" + s[3],
		}, apis.ErrWrongParamsList
	}

	if s := regDublicated.FindStringSubmatch(msg); len(s) > 0 {
		logs.DebugLog("%#v %[1]T", errors.Cause(err))
		return map[string]string{
			s[1]: "duplicate key value violates unique constraint",
		}, apis.ErrWrongParamsList
	}

	return nil, err
}
