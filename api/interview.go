// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"strconv"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/httpgo/services"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

type DTOSendCV struct {
	CompId                  int             `json:"comp_id"`
	FreelancerId            interface{}     `json:"freelancerId"`
	CheckedVacanciesEntries [][]interface{} `json:"checkedVacanciesEntries"`
	CheckedEmailsEntries    [][]interface{} `json:"checkedEmailsEntries"`
	EmailSubject            string          `json:"emailSubject"`
	EmailTemplate           string          `json:"emailTemplate"`
}

func (d *DTOSendCV) GetValue() interface{} {
	return d
}

func (d *DTOSendCV) NewValue() interface{} {
	return &DTOSendCV{}
}
func HandleSendCV(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOSendCV)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	user := auth.GetUserData(ctx)
	timeNow := time.Now()
	tableVTC, _ := db.NewVacancies_to_candidates(DB)
	candidates, _ := db.NewCandidates(DB)
	IntRevCandidate, _ := db.NewInt_rev_candidates(DB)
	SendedEmail, _ := db.NewSended_emails(DB)
	for _, val := range u.CheckedVacanciesEntries {
		s, ok := val[1].(bool)
		if ok && s {
			v, ok := val[0].(string)
			if !ok {
				return "wrong DTO", apis.ErrWrongParamsList
			}

			vacID, err := strconv.Atoi(v)
			if err != nil {
				return "wrong id vacancy", apis.ErrWrongParamsList
			}

			_, err = tableVTC.Upsert(ctx,
				dbEngine.ColumnsForSelect("company_id", "candidate_id", "vacancy_id",
					"status", "user_id", "date_last_change"),
				dbEngine.ArgsForSelect(u.CompId, id, vacID, 9, user.Id, timeNow),
				dbEngine.InsertOnConflict("candidate_id, vacancy_id"),
			)
			if err != nil {
				return createErrResult(err)
			}

			_, err = IntRevCandidate.Insert(ctx,
				dbEngine.ColumnsForSelect("company_id", "candidate_id", "vacancy_id",
					"status", "user_id", "date"),
				dbEngine.ArgsForSelect(u.CompId, id, vacID, 9, user.Id, timeNow),
			)
			if err != nil {
				return createErrResult(err)
			}

			_, err = SendedEmail.Insert(ctx,
				dbEngine.ColumnsForSelect("company_id", "user_id",
					"emails", "subject", "text_emails", "meet_id"),
				dbEngine.ArgsForSelect(u.CompId, user.Id, user.Email, u.EmailSubject, u.EmailTemplate, 0),
			)
			if err != nil {
				return createErrResult(err)
			}

			_, err = candidates.Update(ctx,
				dbEngine.ColumnsForSelect("date"),
				dbEngine.WhereForSelect("id"),
				dbEngine.ArgsForSelect(timeNow, id),
			)
			if err != nil {
				return createErrResult(err)
			}

			toLogCandidateVacancy(ctx, DB, id, int32(u.CompId), int32(vacID), " отправил CV кандидата  ", CODE_LOG_UPDATE)
		}
	}

	for _, val := range u.CheckedEmailsEntries {
		s, ok := val[1].(bool)
		if ok && s {
			email, ok := val[0].(string)
			if !ok {
				return "wrong DTO", apis.ErrWrongParamsList
			}
			err := services.Send(ctx, "mail", services.Mail{
				From:        "cv@uppeople.co",
				To:          email,
				Subject:     u.EmailSubject,
				ContentType: "text/html",
				Body:        u.EmailTemplate,
				Attachments: nil,
			})
			if err != nil {
				return createErrResult(err)
			}
		}
	}

	return u, nil
}
