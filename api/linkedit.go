// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func HandleAuthLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	v, err := HandleAuthLogin(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	m := map[string]interface{}{
		"status":  "success",
		"user_id": v.(map[string]interface{})["user"].(map[string]interface{})["id"],
	}

	return m, nil
}

func HandleGetPlatformsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_platforms",
	}

	return m, nil
}

func HandleGetCandidate_infoLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_candidate_info",
	}

	return m, nil
}

func HandleGetReasonsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_reasons",
	}

	return m, nil
}

func HandleGetRecruiterVacanciesLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_recruiter_vacancies",
	}

	return m, nil
}

func HandleGetTagsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_tags",
	}

	return m, nil
}

func HandleGetSenioritiesLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_seniorities",
	}

	return m, nil
}
func HandleGetSelectorsLinkEdit(ctx *fasthttp.RequestCtx) (interface{}, error) {
	// getPlatforms(ctx, )
	m := map[string]interface{}{
		"status": "get_selectors",
	}

	return m, nil
}
