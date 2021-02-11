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

type ViewCompany struct {
	*db.CompaniesFields
	Vacancies  int32                `json:"vacancies,omitempty"`
	Candidates int32                `json:"candidates,omitempty"`
	Calendar   []*db.MeetingsFields `json:"calendar,omitempty"`
	Contacts   []*db.ContactsFields `json:"contacts,omitempty"`
	Managers   []*db.UsersFields    `json:"managers,omitempty"`
}

func HandleInformationForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	companies, _ := db.NewCompanies(DB)
	err := companies.SelectOneAndScan(ctx,
		companies,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}
	v := &ViewCompany{
		CompaniesFields: companies.Record,
		Calendar:        make([]*db.MeetingsFields, 0),
		Contacts:        make([]*db.ContactsFields, 0),
		Managers:        make([]*db.UsersFields, 0),
	}

	contacts, _ := db.NewContacts(DB)
	err = contacts.SelectSelfScanEach(ctx,
		func(record *db.ContactsFields) error {
			v.Contacts = append(v.Contacts, record)
			return nil
		},
		dbEngine.WhereForSelect("company_id"),
		dbEngine.ArgsForSelect(companies.Record.Id),
	)
	if err != nil {
		logs.ErrorLog(err, "contacts.SelectSelfScanEach")
	}

	meeting, _ := db.NewMeetings(DB)
	err = meeting.SelectSelfScanEach(ctx,
		func(record *db.MeetingsFields) error {
			v.Calendar = append(v.Calendar, record)
			return nil
		},
		dbEngine.WhereForSelect("company_id"),
		dbEngine.ArgsForSelect(companies.Record.Id),
	)
	if err != nil {
		logs.ErrorLog(err, "meeting.SelectSelfScanEach")
	}

	users, _ := db.NewUsers(DB)
	err = users.SelectSelfScanEach(ctx,
		func(record *db.UsersFields) error {
			v.Managers = append(v.Managers, record)
			return nil
		},
		dbEngine.WhereForSelect("<"),
		dbEngine.ArgsForSelect(4),
	)
	if err != nil {
		logs.ErrorLog(err, "users.SelectSelfScanEach")
	}

	return v, nil
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
		dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	}

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
				"select count(*) from vacancies where company_id=$1",
				record.Id,
			)
			if err != nil {
				return errors.Wrap(err, "")
			}
			err = DB.Conn.SelectOneAndScan(ctx,
				&elem.Candidates,
				"select count(candidate_id) from vacancies_to_candidates where company_id=$1",
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
