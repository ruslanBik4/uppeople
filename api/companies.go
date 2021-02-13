// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type SearchCompany struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	IsActive       bool   `json:"IsActive"`
	Skype          string `json:"skype"`
	Phone          string `json:"phone"`
	WithRecruiters bool   `json:"WithRecruiters"`
}

func (s *SearchCompany) GetValue() interface{} {
	return s
}

func (s *SearchCompany) NewValue() interface{} {
	return &SearchCompany{}
}

func HandleAllCompanies(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	offset := 0
	id, ok := ctx.UserValue(ParamPageNum.Name).(int)
	if ok && id > 1 {
		offset = id * pageItem
	}

	companies, _ := db.NewCompanies(DB)
	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy("name"),
		// dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	}

	sqlVacancy := "select count(*) from vacancies where company_id=$1"
	sqlCandidates := "select count(distinct candidate_id) from candidates_to_companies where company_id=$1"
	dto, ok := ctx.UserValue(apis.JSONParams).(*SearchCompany)
	if ok {
		args := make([]interface{}, 0)
		where := make([]string, 0)
		if dto.Name > "" {
			where = append(where, "~name")
			args = append(args, dto.Name)
		}
		if dto.Email > "" {
			where = append(where, "~email")
			args = append(args, dto.Email)
		}
		if dto.Skype > "" {
			where = append(where, "~skype")
			args = append(args, dto.Skype)
		}
		if dto.Phone > "" {
			where = append(where, "~phone")
			args = append(args, dto.Phone)
		}
		if dto.IsActive {
			where = append(where, `id in (SELECT company_id
			FROM vacancies
			WHERE status=%s)`)
			args = append(args, []int32{0, 1})
			fActive := " AND status=ANY(array[0,1])"
			sqlVacancy += fActive
			// sqlCandidates += fActive
		}
		options = append(options,
			dbEngine.WhereForSelect(where...),
			dbEngine.ArgsForSelect(args...),
		)
	} else {
		logs.DebugLog("not json")
	}

	rows := make([]*ViewCompany, 0)
	err := companies.SelectSelfScanEach(ctx,
		func(record *db.CompaniesFields) error {

			elem := &ViewCompany{
				CompaniesFields: record,
			}
			err := DB.Conn.SelectOneAndScan(ctx,
				&elem.Vacancies,
				sqlVacancy,
				record.Id,
			)
			if err != nil {
				return errors.Wrap(err, "")
			}
			err = DB.Conn.SelectOneAndScan(ctx,
				&elem.Candidates,
				sqlCandidates,
				record.Id,
			)
			if err != nil {
				return errors.Wrap(err, "")
			}

			rows = append(rows, elem)

			return nil
		},
		options...)
	if err != nil {
		return createErrResult(err)
	}

	return rows, nil
}
