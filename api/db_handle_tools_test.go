// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/uppeople/db"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestEmptyValue(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  bool
	}{
		{
			"int32",
			int32(0),
			true,
		},
		{
			"int64",
			int64(0),
			true,
		},
		{
			"float32",
			float32(0),
			true,
		},
		{
			"float64",
			float64(0),
			true,
		},
		{
			"sql.NullInt32",
			sql.NullInt32{},
			true,
		},
		{
			"sql.NullInt64",
			sql.NullInt64{},
			true,
		},
		{
			"NullFloat64",
			sql.NullFloat64{},
			true,
		},
		{
			"NullString",
			sql.NullString{},
			true,
		},
		{
			"NullTime",
			sql.NullTime{},
			true,
		},
		// positive case
		{
			"int32",
			int32(1),
			false,
		},
		{
			"int64",
			int64(10),
			false,
		},
		{
			"float32",
			float32(1),
			false,
		},
		{
			"float64",
			float64(-1),
			false,
		},
		{
			"sql.NullInt32",
			sql.NullInt32{0, true},
			false,
		},
		{
			"sql.NullInt64",
			sql.NullInt64{0, true},
			false,
		},
		{
			"NullFloat64",
			sql.NullFloat64{0, true},
			false,
		},
		{
			"NullString",
			sql.NullString{"test", true},
			false,
		},
		{
			"NullTime",
			sql.NullTime{time.Now(), true},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, EmptyValue(tt.value))
		})
	}
}

func Test_getCompanies(t *testing.T) {
	type args struct {
		ctx *fasthttp.RequestCtx
		DB  *dbEngine.DB
		opt []dbEngine.BuildSqlOptions
	}
	tests := []struct {
		name    string
		args    args
		wantRes db.SelectedUnits
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := getCompanies(tt.args.ctx, tt.args.DB, tt.args.opt...); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("getCompanies() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_getRecruiters(t *testing.T) {
	type args struct {
		ctx *fasthttp.RequestCtx
		DB  *dbEngine.DB
	}
	tests := []struct {
		name    string
		args    args
		wantRes db.SelectedUnits
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := getRecruiters(tt.args.ctx, tt.args.DB); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("getRecruiters() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_getVacToCand(t *testing.T) {
	type args struct {
		ctx *fasthttp.RequestCtx
		DB  *dbEngine.DB
	}
	tests := []struct {
		name    string
		args    args
		wantRes db.SelectedUnits
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := getVacToCand(tt.args.ctx, tt.args.DB); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("getVacToCand() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
