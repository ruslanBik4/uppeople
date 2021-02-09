// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
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
	SelectPlatform    SelectedUnit   `json:"selectPlatform"`
	SelectSeniority   SelectedUnit   `json:"selectSeniority"`
	SelectedTag       SelectedUnit   `json:"selectedTag"`
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
	Status   statusCandidate `json:"status"`
}

type VacanciesDTO struct {
	*db.VacanciesFields
	Platforms *db.PlatformsFields `json:"platforms"`
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
	Status_vac       *db.Status_for_vacsFields `json:"status_vac"`
	User_id          int32                     `json:"user_id"`
	Vacancy          VacanciesDTO              `json:"vacancy"`
	Vacancy_id       int32                     `json:"vacancy_id"`
}
type ViewCandidates struct {
	Candidate *ViewCandidate      `json:"0"`
	SelectOpt selectOpt           `json:"select"`
	Statuses  []StatusesCandidate `json:"statusesCandidate"`
}

const pageItem = 15

func HandleViewInformationForSendCV(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
		return nil, errors.Wrap(err, "	")
	}

	return map[string]interface{}{
		"companies": nil,
		"subject":   nil,
		"candId":    id,
		"emailTemplay": map[string]string{
			"text": fmt.Sprintf(`<p><span style="font-size: 14px;">Please, review {platform} %s  CV</span></p>
<p>%s</p>
<p><br>Will be appreciate for quick feedback.</p>
<p><br><br></p>
<p>@"UPpeople" Recruiting agency</p>
<p>&nbsp;<a href="http://www.rock-it.com.ua/" target="_self"><span style="color: blue;font-size: 16px;font-family: Journal, serif;">http://www.rock-it.com.ua/</span></a><span style="font-size: 16px;"> </span></p>`,
				table.Record.Name, table.Record.Link),
		},
	}, nil
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

	return DB.Conn.SelectToMaps(ctx,
		"select * from comments_for_candidates where id=$1 order by created_at DESC",
		id,
	)
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

	maps := make(map[string]interface{}, 0)
	maps["companies"], err = DB.Conn.SelectToMaps(ctx,
		`SELECT id as compId, c.name, otpravka as send_details,
  (select json_agg(json_build_object('email', t.email, 'id',t.id, t.all_platforms,
           'p',contacts_to_platforms.platform_id, 'name',t.name))
             from contacts t join contacts_to_platforms on t.id=contacts_to_platforms.contact_id
           WHERE t.company_id = c.id AND (platform_id=$1 OR all_platforms=1)) as contacts,
  (select json_agg(json_build_object('id', v.id,
                                           'platform', (select p.nazva  from platforms p where p.id = v.platform_id),
                                           'location',
           (select l.name   from location_for_vacancies l where v.location_id = l.id),
                                           'seniority', (select s.nazva   from seniorities s where s.id = v.seniority_id),
           'salary', v.salary, 'name', v.name, 'user_ids', v.user_ids)) as vacancy
	from vacancies v
	where c.id = v.company_id and status <= 1 and platform_id=$1)
FROM companies c
WHERE c.id in (select v.company_id from vacancies v
    where status <= 1 and platform_id=$1)`,
		table.Record.Platform_id)
	if err != nil {
		return createErrResult(err)
	}

	for _, val := range maps["companies"].([]map[string]interface{})[0] {
		logs.DebugLog("%T %#[1]v", val)
	}

	maps["subject"] = fmt.Sprintf("%s UPpeople CV %s - %s", time.Now().Format("02-01-2006"), platformName, name)
	maps["emailTemplay"] = fmt.Sprintf(emailText, platformName, name, table.Record.Link)

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
		return nil, errors.Wrap(err, "	")
	}

	auth.PutEditCandidate(ctx, table.Record)
	res := ViewCandidates{
		SelectOpt: NewSelectOpt(ctx, DB),
		Statuses: []StatusesCandidate{
			{
				Candidate_id: table.Record.Id,
				Company:      &db.CompaniesFields{},
				Status_vac:   &db.Status_for_vacsFields{},
				Vacancy: VacanciesDTO{
					&db.VacanciesFields{},
					&db.PlatformsFields{},
				},
			},
		},
	}

	res.Candidate = NewCandidateView(ctx, table.Record, DB, res.SelectOpt.Platforms, res.SelectOpt.Seniorities).ViewCandidate

	return res, nil
}

func HandleAddCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*CandidateDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	logs.DebugLog(u)
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

	if u.SelectedTag.Id > 0 {
		u.Tag_id = u.SelectedTag.Id
	}
	if u.Tag_id == 0 {
		return map[string]interface{}{
			"tag_id": "required value",
		}, apis.ErrWrongParamsList
	}
	if u.SelectPlatform.Id > 0 {
		u.Platform_id = u.SelectPlatform.Id
	}
	if u.SelectSeniority.Id > 0 {
		u.Seniority_id = u.SelectSeniority.Id
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
		"i":       i,
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
		return nil, err
	}

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

	if u.SelectedTag.Id > 0 {
		u.Tag_id = u.SelectedTag.Id
	}
	if u.Tag_id == 0 {
		return map[string]interface{}{
			"tag_id": "required value",
		}, apis.ErrWrongParamsList
	}
	if u.SelectPlatform.Id > 0 {
		u.Platform_id = u.SelectPlatform.Id
	}
	if u.SelectSeniority.Id > 0 {
		u.Seniority_id = u.SelectSeniority.Id
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
			u.SelectPlatform.Id,
			u.Salary,
			u.Email,
			u.Phone,
			u.Skype,
			u.Link,
			u.Linkedin,
			u.Str_companies,
			u.Status,
			u.SelectedTag.Id,
			u.Comment,
			u.Resume,
			u.Sfera,
			u.Experience,
			u.Education,
			u.Language,
			u.Zapoln_profile,
			u.File,
			// u.Avatar,
			u.SelectSeniority.Id,
			u.Date_follow_up,
		}
	}

	args = append(args, id)

	if u.Tag_id == 3 || u.Tag_id == 4 {
		columns = append(columns, "recruter_id")
		args = append(args, auth.GetUserData(ctx).Id)
	}
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
