// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package data

import (
	"go/types"

	"github.com/ruslanBik4/httpgo/apis"
)

var (
	paramsID = apis.InParam{
		Name: "id",
		Req:  false,
		Type: apis.NewTypeInParam(types.Int32),
	}
	ParamsLang = apis.InParam{
		Name:     "lang",
		Desc:     "need to get result on non-english",
		DefValue: "en",
		Req:      true,
		Type:     apis.NewTypeInParam(types.String),
	}
	ParamsGetFormActions = apis.InParam{
		Name: "is_get_form_actions",
		Desc: "need to get form actions in response",
		Req:  false,
		Type: apis.NewTypeInParam(types.Bool),
	}
	ParamsCounter = apis.InParam{
		Name:     "counter",
		Desc:     "max count of results",
		DefValue: 10,
		Req:      true,
		Type:     apis.NewTypeInParam(types.Int),
	}
	ParamsHTML = apis.InParam{
		Name: "html",
		Desc: "need for get result in html instead JSON",
		Req:  false,
		Type: apis.NewTypeInParam(types.Bool),
	}
	ParamsValue = apis.InParam{
		Name: "value",
		Desc: "part of polymers name, its Synonyms or Abbreviations",
		Req:  true,
		Type: apis.NewTypeInParam(types.String),
	}
)

const PathVersion = "/api/v1"
