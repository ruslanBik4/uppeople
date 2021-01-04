// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

type SelectedUnit struct {
	Id    int32  `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type CandidateDTO struct {
	*db.CandidatesFields
	Comment           string         `json:"comment"`
	Date              string         `json:"date"`
	Phone             string         `json:"phone"`
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
	Platform  *db.PlatformsFields `json:"platforms,omitempty"`
	Seniority string              `json:"seniority"`
	Tags      *db.TagsFields      `json:"tags,omitempty"`
	Recruiter string              `json:"recruiter"`
}

type CandidateView struct {
	*ViewCandidate
	Platform string          `json:"platform,omitempty"`
	TagName  string          `json:"tag_name,omitempty"`
	TagColor string          `json:"tag_color,omitempty"`
	Status   statusCandidate `json:"status"`
}
type ResCandidates struct {
	Count, Page int
	CurrentPage int                     `json:"currentPage"`
	PerPage     int                     `json:"perPage"`
	Candidates  []*CandidateView        `json:"candidates"`
	Company     []*db.CompaniesFields   `json:"company"`
	Platforms   []*db.PlatformsFields   `json:"platforms"`
	Recruiter   []*db.UsersFields       `json:"recruiter"`
	Reasons     []*db.TagsFields        `json:"reasons"`
	Seniority   []*db.SenioritiesFields `json:"seniority"`
	Statuses    []*db.StatusesFields    `json:"statuses"`
	Tags        []*db.TagsFields        `json:"tags"`
}

type selectOpt struct {
	Companies     []*db.CompaniesFields               `json:"companies"`
	Platforms     []*db.PlatformsFields               `json:"platforms"`
	Recruiters    []*db.UsersFields                   `json:"recruiters"`
	Statuses      []*db.StatusesFields                `json:"candidateStatus"`
	Location      []*db.Location_for_vacanciesFields  `json:"location"`
	RejectReasons []*db.TagsFields                    `json:"reject_reasons"`
	RejectTag     []*db.TagsFields                    `json:"reject_tag"`
	Recruiter     []*db.UsersFields                   `json:"recruiter"`
	Seniorities   []*db.SenioritiesFields             `json:"seniorities"`
	Tags          []*db.TagsFields                    `json:"tags"`
	VacancyStatus []*db.Vacancies_to_candidatesFields `json:"vacancyStatus"`
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

func HandleAllCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	candidates, _ := db.NewCandidates(DB)
	res := ResCandidates{
		Page:        1,
		CurrentPage: 1,
		PerPage:     pageItem,
		Candidates:  make([]*CandidateView, pageItem),
		Company:     getCompanies(ctx, DB),
		Reasons:     make([]*db.TagsFields, 0),
		Platforms:   getPlatforms(ctx, DB),
		Recruiter:   getRecruter(ctx, DB),
		Seniority:   getSeniorities(ctx, DB),
		Statuses:    getStatUses(ctx, DB),
		Tags:        getTags(ctx, DB),
	}

	for _, tag := range res.Tags {
		if tag.Parent_id == 3 {
			res.Reasons = append(res.Reasons, tag)
		}
	}

	seniors := getSeniorities(ctx, DB)
	recTable, _ := db.NewUsers(DB)
	i := 0
	err := candidates.SelectSelfScanEach(ctx,
		func(record *db.CandidatesFields) error {

			cV := NewCandidateView(ctx, record, recTable, res.Tags, res.Platforms, seniors)

			res.Candidates[i] = cV
			i++
			if i == pageItem {
				return errLimit
			}

			return nil
		},
		dbEngine.OrderBy("date desc"),
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	return res, nil
}

func NewCandidateView(ctx *fasthttp.RequestCtx,
	record *db.CandidatesFields,
	recTable *db.Users,
	tags []*db.TagsFields,
	platforms []*db.PlatformsFields,
	seniors []*db.SenioritiesFields,
) *CandidateView {
	ref := &CandidateView{
		ViewCandidate: &ViewCandidate{CandidatesFields: record},
		Status: statusCandidate{
			Date:         record.Date,
			Comments:     record.Comments,
			CompId:       0,
			Recruiter:    "",
			DateFollowUp: record.Date_follow_up,
		},
	}
	if record.Recruter_id.Valid {
		err := recTable.SelectOneAndScan(ctx,
			&ref.Recruiter,
			dbEngine.ColumnsForSelect("name"),
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(record.Recruter_id.Int64),
		)
		if err != nil {
			logs.ErrorLog(err, "recTable.SelectOneAndScan")
		}
		ref.Status.Recruiter = ref.Recruiter
	}

	for _, tag := range tags {
		if tag.Id == record.Tag_id {
			ref.TagName = tag.Name
			ref.TagColor = tag.Color
			ref.Tags = tag
		}
	}
	for _, s := range seniors {
		if s.Id == ref.Seniority_id.Int64 {
			ref.Seniority = s.Nazva.String
		}
	}

	for _, p := range platforms {
		if p.Id == record.Platform_id.Int64 {
			ref.Platform = p.Nazva.String
			ref.ViewCandidate.Platform = p
		}
	}

	return ref
}

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
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	tags := getTags(ctx, DB)
	platforms := getPlatforms(ctx, DB)
	seniorities := getSeniorities(ctx, DB)
	res := ViewCandidates{
		SelectOpt: selectOpt{
			Companies:     getCompanies(ctx, DB),
			Platforms:     getPlatforms(ctx, DB),
			Recruiters:    getRecruter(ctx, DB),
			Statuses:      getStatUses(ctx, DB),
			Location:      getLocations(ctx, DB),
			RejectReasons: make([]*db.TagsFields, 0),
			RejectTag:     make([]*db.TagsFields, 0),
			Seniorities:   seniorities,
			Tags:          make([]*db.TagsFields, 0),
			VacancyStatus: getVacToCand(ctx, DB),
		},
		Statuses: []StatusesCandidate{
			{
				Candidate_id: int32(table.Record.Id),
				Company:      &db.CompaniesFields{},
				Status_vac:   &db.Status_for_vacsFields{},
				Vacancy: VacanciesDTO{
					&db.VacanciesFields{},
					&db.PlatformsFields{},
				},
			},
		},
	}

	for _, tag := range tags {
		if tag.Parent_id == 3 {
			res.SelectOpt.RejectReasons = append(res.SelectOpt.RejectReasons, tag)
		}

	}

	recTable, _ := db.NewUsers(DB)

	res.Candidates = NewCandidateView(ctx, table.Record, recTable, tags, platforms, seniorities).ViewCandidate

	return res, nil
}

func getStatUses(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.StatusesFields {
	statUses, _ := db.NewStatuses(DB)
	res := make([]*db.StatusesFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.StatusesFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getVacToCand(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.Vacancies_to_candidatesFields {
	vacCand, _ := db.NewVacancies_to_candidates(DB)
	res := make([]*db.Vacancies_to_candidatesFields, 0)

	err := vacCand.SelectSelfScanEach(ctx,
		func(record *db.Vacancies_to_candidatesFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getStatusVac(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.Status_for_vacsFields {
	statUses, _ := db.NewStatus_for_vacs(DB)
	res := make([]*db.Status_for_vacsFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.Status_for_vacsFields) error {
			res = append(res, record)

			return nil
		},
		dbEngine.OrderBy("order_num"),
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getTags(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.TagsFields {
	statUses, _ := db.NewTags(DB)
	res := make([]*db.TagsFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.TagsFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getLocations(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.Location_for_vacanciesFields {
	statUses, _ := db.NewLocation_for_vacancies(DB)
	res := make([]*db.Location_for_vacanciesFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.Location_for_vacanciesFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getSeniorities(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.SenioritiesFields {
	statUses, _ := db.NewSeniorities(DB)
	res := make([]*db.SenioritiesFields, 0)

	err := statUses.SelectSelfScanEach(ctx,
		func(record *db.SenioritiesFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getPlatforms(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.PlatformsFields {
	platforms, _ := db.NewPlatforms(DB)
	res := make([]*db.PlatformsFields, 0)

	err := platforms.SelectSelfScanEach(ctx,
		func(record *db.PlatformsFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getRecruter(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.UsersFields {
	users, _ := db.NewUsers(DB)
	res := make([]*db.UsersFields, 0)

	err := users.SelectSelfScanEach(ctx,
		func(record *db.UsersFields) error {
			res = append(res, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	")
	}

	return res
}

func getCompanies(ctx *fasthttp.RequestCtx, DB *dbEngine.DB) []*db.CompaniesFields {
	company, _ := db.NewCompanies(DB)
	companies := make([]*db.CompaniesFields, 0)

	err := company.SelectSelfScanEach(ctx,
		func(record *db.CompaniesFields) error {
			companies = append(companies, record)

			return nil
		},
	)
	if err != nil {
		logs.ErrorLog(err, "	SelectSelfScanEach")
	}

	return companies
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

	table, _ := db.NewCandidates(DB)
	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(
			"name",
			"platform_id",
			"salary",
			"email",
			"mobile",
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
			"avatar",
			"seniority_id",
			"date_follow_up",
		),
		dbEngine.ArgsForSelect(
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
			u.Date,
			auth.GetUserData(ctx).Id,
			u.Resume,
			u.Sfera,
			u.Experience,
			u.Education,
			u.Language,
			u.Zapoln_profile,
			u.File,
			u.Avatar,
			u.SelectSeniority.Id,
			u.Date_follow_up,
		),
	)
	if err != nil {
		return createErrResult(err)
	}

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

	log, _ := db.NewLogs(DB)
	user := auth.GetUserData(ctx)
	_, err = log.Insert(ctx,
		dbEngine.ColumnsForSelect("user_id", "candidate_id", "text", "date_create", "d_c",
			"kod_deystviya"),
		dbEngine.ArgsForSelect(user.Id, u.CandidateId,
			fmt.Sprintf("Пользователь %s проработал кандидата #%d. Follow-Up: %v . Comment: %s",
				user.Name, u.CandidateId, u.DateFollowUp, u.Comment),
			time.Now(),
			time.Now(),
			102),
	)
	if err != nil {
		return createErrResult(err)
	}

	return createResult(i)
}

func createResult(i int64) (interface{}, error) {
	return map[string]interface{}{
		"message": "Successfully",
		"i":       i,
	}, nil
}

func HandleEditCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*CandidateDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	p, ok := ctx.UserValue(apis.ChildRoutePath).(string)
	if !ok {
		return "wrong id", apis.ErrWrongParamsList
	}

	id, err := strconv.Atoi(p)
	if err != nil {
		return err.Error(), apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewCandidates(DB)
	columns := dbEngine.ColumnsForSelect(
		"name",
		"platform_id",
		"salary",
		"email",
		"mobile",
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
		"avatar",
		"seniority_id",
		"date_follow_up",
	)
	args := dbEngine.ArgsForSelect(
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
		u.Date,
		auth.GetUserData(ctx).Id,
		u.Resume,
		u.Sfera,
		u.Experience,
		u.Education,
		u.Language,
		u.Zapoln_profile,
		u.File,
		u.Avatar,
		u.SelectSeniority.Id,
		u.Date_follow_up,
		id,
	)
	i, err := table.Update(ctx,
		columns,
		dbEngine.WhereForSelect("id"),
		args,
	)
	if err != nil {
		return createErrResult(err)
	}

	return createResult(i)
}
