// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/httpgo/services"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

type CandidateDTO struct {
	*db.CandidatesFields
	Comment           string         `json:"comment"`
	Date              string         `json:"date"`
	Resume            string         `json:"resume"`
	SelectedVacancies []SelectedUnit `json:"selectedVacancies"`
}

func (c *CandidateDTO) GetValue() interface{} {
	return c
}

func (c *CandidateDTO) NewValue() interface{} {
	return &CandidateDTO{CandidatesFields: &db.CandidatesFields{}}
}

type statusCandidate struct {
	Date         time.Time  `json:"date"`
	Comments     string     `json:"comments"`
	CompId       int32      `json:"comp_id"`
	Recruiter    string     `json:"recruiter"`
	DateFollowUp *time.Time `json:"date_follow_up"`
	VacStat      string     `json:"vacStat"`
	CompName     string     `json:"compName"`
	CommentVac   string     `json:"commentVac"`
}
type ViewCandidate struct {
	*db.CandidatesFields
	Platform  *SelectedUnit            `json:"platforms,omitempty"`
	Companies SelectedUnits            `json:"companies,omitempty"`
	Seniority string                   `json:"seniority"`
	Tags      *db.TagsFields           `json:"tags,omitempty"`
	Recruiter string                   `json:"recruiter"`
	Vacancies []map[string]interface{} `json:"vacancies"`
}

type CandidateView struct {
	*ViewCandidate
	Platform string          `json:"platform,omitempty"`
	TagName  string          `json:"tag_name,omitempty"`
	TagColor string          `json:"tag_color,omitempty"`
	Color    string          `json:"color,omitempty"`
	Status   statusCandidate `json:"status"`
}

type VacanciesDTO struct {
	*db.VacanciesFields
	Platforms      *db.PlatformsFields `json:"platforms"`
	DateLastChange time.Time           `json:"date_last_change"`
}
type StatusesCandidate struct {
	Candidate_id     int32                     `json:"candidate_id"`
	Company          *db.CompaniesFields       `json:"company"`
	Company_id       int32                     `json:"company_id"`
	Date_create      time.Time                 `json:"date_create"`
	Date_last_change time.Time                 `json:"date_last_change"`
	Id               int32                     `json:"id"`
	Notice           string                    `json:"notice"`
	Rating           string                    `json:"rating"`
	Rej_text         string                    `json:"rej_text"`
	Status           int32                     `json:"status"`
	Status_vac       *db.Status_for_vacsFields `json:"vacancyStatus"`
	User_id          int32                     `json:"user_id"`
	Vacancy          VacanciesDTO              `json:"vacancy"`
	Vacancy_id       int32                     `json:"vacancy_id"`
}
type ViewCandidates struct {
	Candidate *ViewCandidate      `json:"0"`
	SelectOpt selectOpt           `json:"select"`
	Statuses  []StatusesCandidate `json:"statuses"`
}

const pageItem = 15

func HandleUpdateStatusCandidates(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*db.Vacancies_to_candidatesFields)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	table, _ := db.NewVacancies_to_candidates(DB)
	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect("status"),
		dbEngine.WhereForSelect("candidate_id", "company_id"),
		dbEngine.ArgsForSelect(u.Status, u.Candidate_id, u.Company_id),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidateVacancy(ctx, DB, u.Candidate_id, u.Company_id, u.Vacancy_id, " изменил статус  кандидата ", CODE_LOG_UPDATE)

	switch u.Status {
	case 2:
		// Meeting::insert(array(
		// 	'candidate_id' => $canridateId,
		// 	'company_id' => $companyId,
		// 	'vacancy_id' => $item->vacancy_id,
		// 	'title' => $title,
		// 	'date' => date('Y-m-d H:i:s'),
		// 	'd_t' => date('Y-m-d'),
		// 	'user_id' => $user->id,
		// 	'color' => null,
		// 	'type' => null
		// 	));
	case 3:
		// IntRevCandidate::insert(array(
		// 	'candidate_id' => $canridateId,
		// 	'company_id' => $companyId,
		// 	'vacancy_id' => $item->vacancy_id,
		// 	'status' => $request->value['id'],
		// 	'user_id' => $user->id,
		// 	'date' => date('Y-m-j')
		// 	));

	}
	return createResult(i)
}

func HandleAddCommentsCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	text := string(ctx.Request.Body())
	table, _ := db.NewComments_for_candidates(DB)
	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect("candidate_id", "comments"),
		dbEngine.ArgsForSelect(id, text),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidate(ctx, DB, id, " оставил комментарий в кандидате "+text, CODE_LOG_UPDATE)

	return createResult(i)

}

func HandleCommentsCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	maps, err := DB.Conn.SelectToMaps(ctx,
		"select * from comments_for_candidates where candidate_id=$1 order by created_at DESC",
		id,
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

func HandleInformationForSendCV(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
	table, _ := db.NewCandidates(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	name := table.Record.Name
	platform, _ := db.NewPlatforms(DB)
	err = platform.SelectOneAndScan(ctx,
		platform,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(table.Record.Platform_id),
	)
	if err != nil {
		return createErrResult(err)
	}

	platformName := platform.Record.Nazva.String

	seniTable, _ := db.NewSeniorities(DB)
	err = seniTable.SelectOneAndScan(ctx,
		seniTable,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(table.Record.Seniority_id),
	)
	if err != nil {
		return createErrResult(err)
	}

	seniority := seniTable.Record.Nazva.String

	maps := make(map[string]interface{}, 0)
	maps["companies"], err = DB.Conn.SelectToMaps(ctx,
		`SELECT id as comp_id, c.name, otpravka as send_details,
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
			'user_ids', v.user_ids)) as vacancy
	from vacancies v
	where c.id = v.company_id and status <= 1 and platform_id=$1)
FROM companies c
WHERE c.id in (select v.company_id from vacancies v
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

	maps["subject"] = fmt.Sprintf("%s UPpeople CV %s - %s", time.Now().Format("02-01-2006"), platformName, name)
	maps["emailTemplay"] = fmt.Sprintf(EMAIL_TEXT, platformName, name, table.Record.Link, seniority,
		table.Record.Language, table.Record.Salary)

	return maps, nil
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
			ParamID.Name: "wrong type, expect int32",
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
	if err != nil && err != pgx.ErrNoRows {
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

	seniTable, _ := db.NewSeniorities(DB)
	err = seniTable.SelectOneAndScan(ctx,
		seniTable,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(record.Seniority_id),
	)
	if err != nil {
		return createErrResult(err)
	}

	seniority := seniTable.Record.Nazva.String

	emailSubject := fmt.Sprintf("UPpeople invite %s - %s", record.Name, platformName)

	emailTemplate := fmt.Sprintf(EMAIL_TEXT, platformName, record.Name, record.Link,
		seniority, record.Language, record.Salary)
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
			ParamID.Name: "wrong type, expect int32",
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
		`SELECT c.id as comp_id, c.name, otpravka as send_details,
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

func HandleViewCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	table, _ := db.NewCandidates(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	auth.PutEditCandidate(ctx, table.Record)
	res := ViewCandidates{
		SelectOpt: NewSelectOpt(ctx, DB),
		Statuses:  []StatusesCandidate{},
	}

	view := NewCandidateView(ctx, table.Record, DB, res.SelectOpt.Platforms, res.SelectOpt.Seniorities)
	res.Candidate = view.ViewCandidate

	for _, vacancy := range res.Candidate.Vacancies {
		res.Statuses = append(res.Statuses, StatusesCandidate{
			Candidate_id: table.Record.Id,
			Company_id:   vacancy["company_id"].(int32),
			Company: &db.CompaniesFields{
				Id: int64(vacancy["company_id"].(int32)),
				Name: sql.NullString{
					String: vacancy["name"].(string),
					Valid:  true,
				},
			},
			Status_vac: &db.Status_for_vacsFields{
				Id: vacancy["status_id"].(int64),
				Status: sql.NullString{
					String: vacancy["status"].(string),
					Valid:  true,
				},
			},
			Vacancy: VacanciesDTO{
				&db.VacanciesFields{
					Id:     vacancy["id"].(int32),
					Salary: vacancy["salary"].(int32),
				},
				&db.PlatformsFields{Nazva: sql.NullString{
					String: vacancy["platform"].(string),
					Valid:  true,
				}},
				vacancy["date_last_change"].(time.Time),
			},
			Date_last_change: vacancy["date_last_change"].(time.Time),
		})
	}
	return res, nil
}

func HandleAddCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*CandidateDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	columns := []string{
		"name",
		"platform_id",
		"salary",
		"email",
		"phone",
		"skype",
		"link",
		"linkedin",
		"str_companies",
		"status",
		"tag_id",
		"comments",
		"date",
		"recruter_id",
		"text_rezume",
		"sfera",
		"experience",
		"education",
		"language",
		"zapoln_profile",
		"file",
		"seniority_id",
		"date_follow_up",
	}

	if u.Tag_id == 0 {
		return map[string]interface{}{
			"tag_id": "required value",
		}, apis.ErrWrongParamsList
	}
	args := []interface{}{
		u.Name,
		u.Platform_id,
		u.Salary,
		u.Email,
		u.Phone,
		u.Skype,
		u.Link,
		u.Linkedin,
		u.Str_companies,
		u.Status,
		u.Tag_id,
		u.Comment,
		time.Now(),
		auth.GetUserData(ctx).Id,
		u.Resume,
		u.Sfera,
		u.Experience,
		u.Education,
		u.Language,
		u.Zapoln_profile,
		u.File,
		u.Seniority_id,
		u.Date_follow_up,
	}

	if u.Avatar > "" {
		columns = append(columns, "avatar")
		args = append(args, u.Avatar)
	}
	table, _ := db.NewCandidates(DB)
	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	err = table.SelectOneAndScan(ctx,
		&u.Id,
		dbEngine.ColumnsForSelect("id"),
		dbEngine.WhereForSelect("name"),
		dbEngine.ArgsForSelect(u.Name),
	)
	if err != nil {
		logs.ErrorLog(err, "table.SelectOneAndScan")
	}
	toLogCandidate(ctx, DB, int32(u.Id), u.Comment, 101)

	ctx.SetStatusCode(fasthttp.StatusCreated)

	return i, nil
}

type FollowUpDTO struct {
	CandidateId  int32  `json:"candidate_id"`
	DateFollowUp string `json:"date_follow_up"`
	Comment      string `json:"comment"`
}

func (f *FollowUpDTO) GetValue() interface{} {
	return f
}

func (f *FollowUpDTO) NewValue() interface{} {
	return &FollowUpDTO{
		// DateFollowUp: &time.Time{},
	}
}

func HandleFollowUpCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*FollowUpDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	candidates, _ := db.NewCandidates(DB)

	i, err := candidates.Update(ctx,
		dbEngine.ColumnsForSelect("date", "date_follow_up", "comments"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(time.Now(), u.DateFollowUp, u.Comment, u.CandidateId),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidate(ctx, DB, u.CandidateId,
		fmt.Sprintf("Follow-Up: %v . Comment: %s", u.DateFollowUp, u.Comment), 102)

	return createResult(i)
}

func createResult(i int64) (interface{}, error) {
	return map[string]interface{}{
		"message": "Successfully",
		"id":      i,
	}, nil
}

func HandleDeleteCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	err := DB.Conn.ExecDDL(ctx, "delete from candidates where id = $1", id)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidate(ctx, DB, id, "", 103)
	toLogCandidate(ctx, DB, id, "", 103)
	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}

func HandleEditCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	u, ok := ctx.UserValue(apis.JSONParams).(*CandidateDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	if u.Tag_id == 0 {
		return map[string]interface{}{
			"tag_id": "required value",
		}, apis.ErrWrongParamsList
	}
	u.Comments = u.Comment

	oldData := auth.GetEditCandidate(ctx)
	table, _ := db.NewCandidates(DB)
	columns := make([]string, 0)
	args := make([]interface{}, 0)
	stopColumns := map[string]bool{
		"recruter_id": true,
		"id":          true,
		"date":        true,
		"avatar":      true,
	}
	if oldData != nil {
		for _, col := range table.Columns() {
			name := col.Name()
			if stopColumns[name] {
				continue
			}

			if oldData.ColValue(name) != u.ColValue(name) {
				columns = append(columns, name)
				args = append(args, u.ColValue(name))
			}
		}
	} else {
		columns = []string{
			"name",
			"platform_id",
			"salary",
			"email",
			"phone",
			"skype",
			"link",
			"linkedin",
			"str_companies",
			"status",
			"tag_id",
			"comments",
			"text_rezume",
			"sfera",
			"experience",
			"education",
			"language",
			"zapoln_profile",
			"file",
			// "avatar",
			"seniority_id",
			"date_follow_up",
		}
		args = []interface{}{
			u.Name,
			u.Platform_id,
			u.Salary,
			u.Email,
			u.Phone,
			u.Skype,
			u.Link,
			u.Linkedin,
			u.Str_companies,
			u.Status,
			u.Tag_id,
			u.Comment,
			u.Resume,
			u.Sfera,
			u.Experience,
			u.Education,
			u.Language,
			u.Zapoln_profile,
			u.File,
			// u.Avatar,
			u.Seniority_id,
			u.Date_follow_up,
		}
	}

	args = append(args, id)

	if u.Tag_id == 3 || u.Tag_id == 4 {
		columns = append(columns, "recruter_id")
		args = append(args, auth.GetUserData(ctx).Id)
	}

	logs.DebugLog(columns, args)
	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	if i > 0 {
		text := ""
		for i, col := range columns {
			if i > 0 {
				text += ", "
			}

			text += fmt.Sprintf("%s=%v", col, args[i])
		}
		toLogCandidate(ctx, DB, int32(id), text, 100)
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}
