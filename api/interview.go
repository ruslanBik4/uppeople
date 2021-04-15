// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/httpgo/services"
	"github.com/ruslanBik4/logs"
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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
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

	table, _ := db.NewCandidates_to_companies(DB)
	_, err := table.Upsert(ctx,
		dbEngine.ColumnsForSelect("company_id", "candidate_id", "visible"),
		dbEngine.ArgsForSelect(u.CompId, id, 0),
	)
	if err != nil {
		return createErrResult(err)
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

	tableCandidate, _ := db.NewCandidates(DB)
	i, err := tableCandidate.Update(ctx,
		dbEngine.ColumnsForSelect("tag_id"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(2, id),
	)
	if err != nil {
		return createErrResult(err)
	}
	if i == 0 {
		logs.DebugLog("Tag_id not updated for candidate id %d on SendCV", id)
	}

	return u, nil
}

type DTOSendInterview struct {
	SelectedCompany  SelectedUnit `json:"selectedCompany"`
	SelectedVacancy  SelectedUnit `json:"selectedVacancy"`
	SelectedContacts []struct {
		Email string `json:"email"`
	} `json:"selectedContacts"`
	Date    string `json:"date"`
	Time    string `json:"time"`
	Comment string `json:"comment"`
}

func (d *DTOSendInterview) GetValue() interface{} {
	return d
}

func (d *DTOSendInterview) NewValue() interface{} {
	return &DTOSendInterview{}
}

func HandleInviteOnInterviewSend(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
	candidates, _ := db.NewCandidates(DB)
	err := candidates.SelectOneAndScan(ctx,
		candidates,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOSendInterview)
	if !ok {
		return fmt.Sprintf("wrong DTO - %T", ctx.UserValue(apis.JSONParams)),
			apis.ErrWrongParamsList
	}

	if len(u.SelectedContacts) < 1 {
		return "need leas a one contact",
			apis.ErrWrongParamsList

	}

	vacID := u.SelectedVacancy.Id
	user := auth.GetUserData(ctx)
	timeNow := time.Now()
	tableVTC, _ := db.NewVacancies_to_candidates(DB)
	IntRevCandidate, _ := db.NewInt_rev_candidates(DB)
	SendedEmail, _ := db.NewSended_emails(DB)

	_, err = tableVTC.Upsert(ctx,
		dbEngine.ColumnsForSelect("company_id", "candidate_id", "vacancy_id",
			"status", "user_id", "date_last_change"),
		dbEngine.ArgsForSelect(u.SelectedCompany.Id, id, vacID, 2, user.Id, timeNow),
	)
	if err != nil {
		return createErrResult(err)
	}

	_, err = IntRevCandidate.Upsert(ctx,
		dbEngine.ColumnsForSelect("company_id", "candidate_id", "vacancy_id",
			"status", "user_id", "date"),
		dbEngine.ArgsForSelect(u.SelectedCompany.Id, id, vacID, 2, user.Id, timeNow),
	)
	if err != nil {
		return createErrResult(err)
	}

	record := candidates.Record
	platform, _ := db.NewPlatforms(DB)
	err = platform.SelectOneAndScan(ctx,
		platform,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(record.Platform_id),
	)
	if err != nil {
		return createErrResult(err)
	}

	platformName := platform.Record.Nazva.String

	emailSubject := fmt.Sprintf("UPpeople invite %s - %s", record.Name, platformName)

	emailTemplate := fmt.Sprintf(EMAIL_INVITE_TEXT, platformName, record.Name,
		u.Date, u.Time, record.Link) + u.Comment
	_, err = SendedEmail.Insert(ctx,
		dbEngine.ColumnsForSelect("company_id", "user_id",
			"emails", "subject", "text_emails", "meet_id"),
		dbEngine.ArgsForSelect(u.SelectedCompany.Id, user.Id, user.Email, emailSubject, emailTemplate, 0),
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

	toLogCandidateVacancy(ctx, DB, id, u.SelectedCompany.Id, vacID, " назначил встречу кандидата ", CODE_LOG_UPDATE)

	tableCTC, _ := db.NewCandidates_to_companies(DB)
	_, err = tableCTC.Upsert(ctx,
		dbEngine.ColumnsForSelect("company_id", "candidate_id", "visible"),
		dbEngine.ArgsForSelect(u.SelectedCompany.Id, id, 0),
	)
	if err != nil {
		return createErrResult(err)
	}

	for _, val := range u.SelectedContacts {
		err := services.Send(ctx, "mail", services.Mail{
			From:        "cv@uppeople.co",
			To:          val.Email,
			Subject:     emailSubject,
			ContentType: "text/html",
			Body:        emailTemplate,
			Attachments: nil,
		})
		if err != nil {
			return createErrResult(err)
		}

	}
	// 	todo: add meeting
	// 	$meetId = Meeting::insertGetId(array(
	// 		'candidate_id' => $request->candidate_id,
	// 		'company_id' => $request->selectedCompany['compId'],
	// 		'vacancy_id' => $request->selectedVacancy['vacId'],
	// 		'title' => $candidate[0]->name,
	// 		'date' => $dataTime,
	// 		'd_t' => $request->date,
	// 		'user_id' => $user->id,
	// ));

	return nil, nil
}

func HandleInviteOnInterviewView(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
	table, _ := db.NewCandidates(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	maps, err := DB.Conn.SelectToMaps(ctx,
		`SELECT c.id as comp_id, c.name, send_details,
  (select json_agg(json_build_object('email', t.email, 'id',t.id, 'all_platforms', t.all_platforms,
           'platform_id', cp.platform_id, 'name', t.name))
             from contacts t left join contacts_to_platforms cp on t.id=cp.contact_id
           WHERE t.company_id = c.id AND (all_platforms=1 OR platform_id=$1)) as contacts,
  (select json_agg(json_build_object('id', v.id,
		   'platform', (select p.nazva  from platforms p where p.id = v.platform_id),
		   'location', (select l.name   from location_for_vacancies l where v.location_id = l.id),
           'seniority', (select s.nazva from seniorities s where s.id = v.seniority_id),
           'salary', v.salary, 
			'name', v.name, 
			'user_ids', v.user_ids)) as vacancies
	from vacancies v
	where c.id = v.company_id and status <= 1 and platform_id=$1)
FROM  companies c 
WHERE exists (select NULL from vacancies_to_candidates where c.id = company_id and status = 9) 
		and
      c.id in (select v.company_id from vacancies v
               where status <= 1 and platform_id=$1 and $2 = ANY(user_ids))`,
		table.Record.Platform_id,
		auth.GetUserData(ctx).Id,
	)
	if err != nil {
		return createErrResult(err)
	}

	if len(maps) == 0 {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return nil, nil
	}

	return maps, nil
}
