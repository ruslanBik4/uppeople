// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"
)

type dtoField map[string]interface{}

func (d *dtoField) CheckParams(ctx *fasthttp.RequestCtx, badParams map[string]string) bool {
	for key, val := range *d {
		if strings.HasSuffix(key, "[]") {
			// key = strings.TrimSuffix(key, "[]")
			switch v := val.(type) {
			case string:
				val = []string{v}
			case []interface{}:
				s := make([]string, len(v))
				for i, str := range v {
					s[i] = fmt.Sprintf("%v", str)
				}
				val = s
			}
		}
		ctx.SetUserValue(key, val)
	}

	return true
}

func (d dtoField) GetValue() interface{} {
	return d
}

func (d dtoField) NewValue() interface{} {
	d = make(map[string]interface{}, 0)

	return &d
}

type FormActions struct {
	Typ string `json:"type"`
	Url string `json:"url"`
}
