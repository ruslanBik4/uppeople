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
}
type ViewCandidate struct {
	*db.CandidatesFields
	Platform  *SelectedUnit            `json:"platforms,omitempty"`
	Companies *SelectedUnit            `json:"companies,omitempty"`
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
	Candidates *ViewCandidate      `json:"0"`
	SelectOpt  selectOpt           `json:"select"`
	Statuses   []StatusesCandidate `json:"statusesCandidate"`
}

const pageItem = 15

func HandleViewCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewCandidates(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(ctx.UserValue("id")),
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

	res.Candidates = NewCandidateView(ctx, table.Record, DB, res.SelectOpt.Platforms, res.SelectOpt.Seniorities).ViewCandidate
	res.Candidates.Vacancies, err = DB.Conn.SelectToMaps(ctx,
		`select vacancies.id, concat(companies.name, ' ("', platforms.nazva, '")') as name, 
LOWER(CONCAT(companies.name, ' ("', platforms.nazva , ')"')) as label, user_ids, platform_id,
		companies, vacancies.company_id, companies.id
FROM vacancies JOIN companies on (vacancies.company_id=companies.id)
	JOIN vacancies_to_candidates on (vacancies.id = vacancies_to_candidates.vacancy_id)
	JOIN platforms ON (vacancies.platform_id = platforms.id)
	WHERE vacancies_to_candidates.candidate_id=$1`, res.Candidates.Id)

	if err != nil {
		logs.ErrorLog(err, "")
	}

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
	args := []interface{}{
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
		time.Now(),
		auth.GetUserData(ctx).Id,
		u.Resume,
		u.Sfera,
		u.Experience,
		u.Education,
		u.Language,
		u.Zapoln_profile,
		u.File,
		u.SelectSeniority.Id,
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

	return createResult(i)
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
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*CandidateDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	u.Platform_id.Int32 = u.SelectPlatform.Id
	u.Platform_id.Valid = true
	u.Seniority_id.Int32 = u.SelectSeniority.Id
	u.Seniority_id.Valid = true
	u.Tag_id = u.SelectedTag.Id
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

	return createResult(i)
}
