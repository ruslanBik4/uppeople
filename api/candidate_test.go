// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"reflect"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

func TestCandidateDTO_GetValue(t *testing.T) {
	type fields struct {
		CandidatesFields *db.CandidatesFields
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CandidateDTO{
				CandidatesFields: tt.fields.CandidatesFields,
			}
			if got := c.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCandidateDTO_NewValue(t *testing.T) {
	type fields struct {
		CandidatesFields *db.CandidatesFields
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CandidateDTO{
				CandidatesFields: tt.fields.CandidatesFields,
			}
			if got := c.NewValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleAddCandidate(t *testing.T) {
	type args struct {
		ctx *fasthttp.RequestCtx
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleAddCandidate(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleAddCandidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleAddCandidate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleAllCandidate(t *testing.T) {
	type args struct {
		ctx *fasthttp.RequestCtx
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleAllCandidate(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleAllCandidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleAllCandidate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
