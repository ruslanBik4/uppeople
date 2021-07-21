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
	"github.com/pkg/errors"
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
	Color        string     `json:"color,omitempty"`
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
	Companies         db.SelectedUnits         `json:"companies,omitempty"`
	Seniority         string                   `json:"seniority"`
	Tags              *db.TagsFields           `json:"tags,omitempty"`
	Recruiter         string                   `json:"recruiter"`
	Vacancies         []map[string]interface{} `json:"vacancies"`
	SelectedVacancies db.SelectedUnits         `json:"selectedVacancies"`
}

type CandidateView struct {
	*ViewCandidate
	TagName  string             `json:"tag_name,omitempty"`
	TagColor string             `json:"tag_color,omitempty"`
	Statuses []*statusCandidate `json:"statuses"`
}

type VacanciesDTO struct {
	*db.VacanciesFields
	Platforms      *db.PlatformsFields `json:"platforms"`
	DateLastChange time.Time           `json:"date_last_change"`
}
type StatusesCandidate struct {
	Candidate_id     int32                   `json:"candidate_id"`
	Company          *db.CompaniesFields     `json:"company"`
	Company_id       int32                   `json:"company_id"`
	Date_create      time.Time               `json:"date_create"`
	Date_last_change time.Time               `json:"date_last_change"`
	Id               int32                   `json:"id"`
	Notice           string                  `json:"notice"`
	Rating           string                  `json:"rating"`
	Rej_text         string                  `json:"rej_text"`
	Status           int32                   `json:"status"`
	Status_vac       *db.StatusForVacsFields `json:"vacancyStatus"`
	User_id          int32                   `json:"user_id"`
	Vacancy          VacanciesDTO            `json:"vacancy"`
	Vacancy_id       int32                   `json:"vacancy_id"`
}
type ViewCandidates struct {
	Candidate *ViewCandidate      `json:"0"`
	SelectOpt selectOpt           `json:"select"`
	Statuses  []StatusesCandidate `json:"statuses"`
}

const pageItem = 15

func HandleReContactCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
	columns := []string{
		"recruter_id",
		"date",
	}
	args := []interface{}{
		auth.GetUserID(ctx),
		time.Now(),
	}

	_, err := table.Update(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(append(args, id)...),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidateRecontact(ctx, id, "")

	return nil, nil
}

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
		dbEngine.WhereForSelect("candidate_id", "vacancy_id"),
		dbEngine.ArgsForSelect(u.Status, u.Candidate_id, u.Vacancy_id),
	)
	if err != nil {
		return createErrResult(err)
	}

	if i == 0 {
		return map[string]string{
			"candidate_id, vacancy_id": "can't find record for this primary",
		}, apis.ErrWrongParamsList
	}

	text := map[string]interface{}{
		"status_for_vac": u.Status,
		"vacancy_id":     u.Vacancy_id,
	}

	if i > 0 {
		toLogCandidateUpdateStatus(ctx, u.Candidate_id, u.Vacancy_id, text)
	}

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

func HandleRmCommentsCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	table, _ := getTableCommentsForCandidates(ctx)
	var returnVal []interface{}
	text := ""
	err := table.SelectOneAndScan(ctx, &returnVal,
		dbEngine.ColumnsForSelect("text_comment", "candidate_id"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id))

	if err != nil {
		return createErrResult(errors.Wrap(err, "comment not found"))
	}

	text = returnVal[0].(string)
	candidateId := returnVal[1].(int32)

	i, err := table.Delete(ctx,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidateDelComment(ctx, candidateId, text)

	return createResult(i)
}

func HandleAddCommentsCandidate(ctx *fasthttp.RequestCtx) (interface{}, error) {
	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	table, _ := getTableCommentsForCandidates(ctx)
	text := string(ctx.Request.Body())
	text = strings.Replace(text, "\"", "", -1)
	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect("candidate_id", "comments"),
		dbEngine.ArgsForSelect(id, text),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogCandidateAddComment(ctx, id, text)

	return createResult(i)
}

func getTableCommentsForCandidates(ctx *fasthttp.RequestCtx) (*db.Comments_for_candidates, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	table, _ := db.NewComments_for_candidates(DB)

	return table, nil
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
		dbEngine.ArgsForSelect(table.Record.Platforms),
	)
	if err != nil {
		return createErrResult(err)
	}

	platformName := platform.Record.Name
	seniority := db.GetSeniorityFromId(table.Record.SeniorityId)

	maps := make(map[string]interface{}, 0)
	maps["companies"], err = DB.Conn.SelectToMaps(ctx,
		SEND_CV_COMPANIES_SQL,
		table.Record.Platforms,
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
	maps["emailTemplay"] = fmt.Sprintf(EMAIL_TEXT, name, platformName, table.Record.Link,
		seniority.Name, "table.Record.IdLanguages", table.Record.Salary)
	// todo add langueages cache
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

	view := NewCandidateView(ctx, table.Record, DB)

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
			Id:         vacancy["id"].(int32),
			Vacancy_id: vacancy["id"].(int32),
			Status_vac: &db.StatusForVacsFields{
				Id:     vacancy["status_id"].(int32),
				Status: vacancy["status"].(string),
			},
			Vacancy: VacanciesDTO{
				&db.VacanciesFields{
					Id:     vacancy["id"].(int32),
					Salary: vacancy["salary"].(int32),
				},
				&db.PlatformsFields{Name: vacancy["platform"].(string)},
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

	table, err := db.NewCandidates(DB)
	if err != nil {
		return createErrResult(err)
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)

	stopColumns := map[string]interface{}{
		"id":          nil,
		"tag_id":      db.GetTagIdFirstContact(),
		"date":        time.Now(),
		"recruter_id": auth.GetUserID(ctx),
	}

	if u != nil {
		for _, col := range table.Columns() {
			name := col.Name()
			if val, ok := stopColumns[name]; ok {
				if val != nil {
					columns = append(columns, name)
					args = append(args, val)
				}
				continue
			}

			newValue := u.ColValue(name)
			if !EmptyValue(newValue) {
				columns = append(columns, name)
				args = append(args, newValue)
			}
		}
	}

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
	toLogCandidateInsert(ctx, u.Id, u.Comments)

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

	toLogCandidatePerform(ctx, u.CandidateId,
		fmt.Sprintf("Follow-Up: %v . Comment: %s", u.DateFollowUp, u.Comment))

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

	toLogCandidateDelete(ctx, id, "")
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
				if name == "vacancies" || name == "platforms" {
					val, ok := oldData.ColValue(name).([]interface{})
					newVal, newOk := newValue.([]interface{})

					if !ok || !newOk {
						continue
					}
					finalValue := make([]interface{}, 0)
					var oldIntVal, newIntVal []int32
					for _, i_val := range val {
						if check_val, ok := i_val.(int32); ok {
							oldIntVal = append(oldIntVal, check_val)
						}
					}

					for _, i_val := range newVal {
						if check_val, ok := i_val.(int32); ok {
							newIntVal = append(newIntVal, check_val)
						}
					}

					for _, ii_val := range newIntVal {
						for _, jj_val := range oldIntVal {
							if ii_val == jj_val {
								continue
							} else {
								finalValue = append(finalValue, ii_val)
								break
							}
						}
					}

					for _, kk_val := range oldIntVal {
						for _, ll_val := range newIntVal {
							if kk_val == ll_val {
								continue
							} else {
								finalValue = append(finalValue, kk_val)
								break
							}
						}
					}

					if len(finalValue) > 0 {
						columns = append(columns, name)
						args = append(args, finalValue)
					}

				}

				columns = append(columns, name)
				args = append(args, newValue)
			}
		}
	} else {
		columns = []string{
			"name",
			"platforms",
			"salary",
			"email",
			"phone",
			"skype",
			"link",
			"linkedin",
			"status",
			"tag_id",
			"comments",
			"cv",
			"experience",
			"education",
			"id_languages",
			"file",
			// "avatar",
			"seniority_id",
			"date_follow_up",
			"vacancies",
		}
		args = []interface{}{
			u.Name,
			u.Platforms,
			u.Salary,
			u.Email,
			u.Phone,
			u.Skype,
			u.Link,
			u.Linkedin,
			u.Status,
			u.Tag_id,
			u.Comments,
			u.Cv,
			u.Experience,
			u.Education,
			u.IdLanguages,
			u.File,
			// u.Avatar,
			u.SeniorityId,
			u.DateFollowUp,
			u.Vacancies,
		}
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
		toLogCandidateUpdate(ctx, id, toLogUpdateValues(columns, args))
		ctx.SetStatusCode(fasthttp.StatusAccepted)
	}

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
