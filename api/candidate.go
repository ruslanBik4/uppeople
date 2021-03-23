// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

// todo - shrink struct
type CandidateDTO struct {
	*db.CandidatesFields
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
	Platform          *SelectedUnit            `json:"platforms,omitempty"`
	Companies         SelectedUnits            `json:"companies,omitempty"`
	Seniority         string                   `json:"seniority"`
	Tags              *db.TagsFields           `json:"tags,omitempty"`
	Recruiter         string                   `json:"recruiter"`
	Vacancies         []map[string]interface{} `json:"vacancies"`
	SelectedVacancies SelectedUnits            `json:"selectedVacancies"`
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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
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
		`SELECT id as comp_id, c.name, send_details,
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
	maps["emailTemplay"] = fmt.Sprintf(EMAIL_TEXT, name, platformName, table.Record.Link, seniority,
		table.Record.Language, table.Record.Salary)

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

	auth.PutEditCandidate(ctx, table.Record)
	res := ViewCandidates{
		SelectOpt: NewSelectOpt(ctx, DB),
		Statuses:  []StatusesCandidate{},
	}

	view := NewCandidateView(ctx, table.Record, DB, res.SelectOpt.Platforms, res.SelectOpt.Seniorities)

	err = DB.Conn.SelectAndScanEach(ctx,
		nil,
		&view.SelectedVacancies,
		`select v.id, 
		concat(companies.name, ' ("', platforms.nazva, '")') as label, 
		LOWER(CONCAT(companies.name, ' ("', platforms.nazva , '")')) as value
	FROM vacancies v JOIN companies on (v.company_id=companies.id)
	JOIN platforms ON (v.platform_id = platforms.id)
	WHERE v.id=ANY($1)
`,
		view.CandidatesFields.Vacancies,
	)
	if err != nil {
		return createErrResult(err)
	}

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
			Id: vacancy["id"].(int32),
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
		"Vacancies",
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
		u.Comments,
		time.Now(),
		auth.GetUserID(ctx),
		u.Text_rezume,
		u.Sfera,
		u.Experience,
		u.Education,
		u.Language,
		u.Zapoln_profile,
		u.File,
		u.Seniority_id,
		u.Date_follow_up,
		u.Vacancies,
	}

	if u.Avatar > "" {
		columns = append(columns, "avatar")
		args = append(args, u.Avatar)
	}
	table, _ := db.NewCandidates(DB)
	id, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	if id <= 0 {
		err = db.LastErr
		if err != (*pgconn.PgError)(nil) {
			return createErrResult(err)
		}

		return DB.Conn.LastRowAffected(), apis.ErrWrongParamsList
	}

	u.Id = int32(id)
	toLogCandidate(ctx, DB, u.Id, u.Comments, CODE_LOG_INSERT)

	ctx.SetStatusCode(fasthttp.StatusCreated)
	//putVacancies(ctx, u, DB)

	return id, nil
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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	err := DB.Conn.ExecDDL(ctx, "delete from candidates where id = $1", id)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidate(ctx, DB, id, "", CODE_LOG_DELETE)
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

	u.Id = id
	if u.Tag_id == 0 {
		return map[string]interface{}{
			"tag_id": "required value",
		}, apis.ErrWrongParamsList
	}

	oldData := auth.GetEditCandidate(ctx)
	table, _ := db.NewCandidates(DB)
	columns := make([]string, 0)
	args := make([]interface{}, 0)
	stopColumns := map[string]bool{
		"recruter_id": true,
		"id":          true,
		"date":        true,
		"avatar":      true,
		"status":      true,
	}
	if oldData != nil {
		for _, col := range table.Columns() {
			name := col.Name()
			if stopColumns[name] {
				continue
			}

			newValue := u.ColValue(name)
			if !EmptyValue(newValue) && (strings.HasPrefix(col.Type(), "_") || oldData.ColValue(name) != newValue) {
				columns = append(columns, name)
				args = append(args, newValue)
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
			"vacancies",
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
			u.Comments,
			u.Text_rezume,
			u.Sfera,
			u.Experience,
			u.Education,
			u.Language,
			u.Zapoln_profile,
			u.File,
			// u.Avatar,
			u.Seniority_id,
			u.Date_follow_up,
			u.Vacancies,
		}
	}

	if u.Tag_id == 3 || u.Tag_id == 4 {
		columns = append(columns, "recruter_id")
		args = append(args, auth.GetUserData(ctx).Id)
	}

	if len(columns) == 0 {
		return "no new data on record", apis.ErrWrongParamsList
	}

	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(append(args, u.Id)...),
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
		toLogCandidate(ctx, DB, id, text, CODE_LOG_UPDATE)
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)
	// putVacancies(ctx, u, DB)

	return nil, nil
}

func putVacancies(ctx *fasthttp.RequestCtx, u *CandidateDTO, DB *dbEngine.DB) {
	if len(u.Vacancies) > 0 {
		table, _ := db.NewVacancies_to_candidates(DB)
		for _, id := range u.Vacancies {

			_, err := table.Upsert(ctx,
				dbEngine.ColumnsForSelect("candidate_id", "vacancy_id", "status"),
				dbEngine.ArgsForSelect(u.Id, id, 1),
			)
			if err != nil {
				logs.ErrorLog(err, "NewVacancies_to_candidates")
			}
		}
	}
}
