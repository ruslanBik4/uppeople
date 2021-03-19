// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"reflect"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

type VacancyDTO struct {
	*db.VacanciesFields
	Comment               string         `json:"comment"`
	SelectCompany         SelectedUnit   `json:"selectCompany"`
	SelectLocation        SelectedUnit   `json:"selectLocation"`
	SelectPlatform        SelectedUnit   `json:"selectPlatform"`
	SelectSeniority       SelectedUnit   `json:"selectSeniority"`
	SelectRecruiter       []SelectedUnit `json:"selectRecruiter"`
	SelectedVacancyStatus int32          `json:"selectedVacancyStatus"`
}

func (v *VacancyDTO) GetValue() interface{} {
	return v
}

func (v *VacancyDTO) NewValue() interface{} {
	return &VacancyDTO{
		VacanciesFields: &db.VacanciesFields{},
	}
}

type vacDTO struct {
	CompanyId             int32          `json:"company_id"`
	Sort                  int32          `json:"sort"`
	CurrentColumn         string         `json:"currentColumn"`
	SelectPlatforms       []SelectedUnit `json:"selectPlatforms"`
	SelectSeniorities     []SelectedUnit `json:"selectSeniorities"`
	SelectCandidateStatus []SelectedUnit `json:"selectCandidate_status"`
	SelectStatuses        []SelectedUnit `json:"selectStatuses"`
}

func (v *vacDTO) GetValue() interface{} {
	return v
}

func (v *vacDTO) NewValue() interface{} {
	return &vacDTO{}
}

func HandleViewVacancy(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	v := &VacanciesView{
		VacanciesFields: &db.VacanciesFields{},
	}
	err := DB.Conn.SelectOneAndScan(ctx,
		v,
		`select *, 
			(select p.nazva from platforms p where v.platform_id=p.id) as platform,
			(select s.nazva from seniorities s where v.seniority_id=s.id) as seniority,
			(select c.name from companies c where v.company_id=c.id) as company,
			(select s.name from location_for_vacancies s where v.location_id=s.id) as location
			from vacancies v
			where id = $1
`,
		ctx.UserValue("id"),
	)
	if err == pgx.ErrNoRows {
		ctx.SetStatusCode(fasthttp.StatusNoContent)
		return nil, nil
	}
	if err != nil {
		return createErrResult(err)
	}

	v.Date = v.Date_create.Format("2006-01-02")

	auth.PutEditVacancy(ctx, v.VacanciesFields)

	return v, nil
}

func HandleEditStatusVacancy(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*VacancyDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	table, _ := db.NewVacancies(DB)
	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect("status"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(u.Status, u.Id),
	)
	if err != nil {
		return createErrResult(err)
	}

	if i > 0 {
		text := fmt.Sprintf("status=%d", u.Id)
		toLogVacancy(ctx, DB, u.SelectCompany.Id, u.Id, text, CODE_LOG_UPDATE)
	}

	return createResult(i)
}

func HandleEditVacancy(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
	u, ok := ctx.UserValue(apis.JSONParams).(*VacancyDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	table, _ := db.NewVacancies(DB)
	oldData := auth.GetEditVacancy(ctx)
	columns := make([]string, 0)
	args := make([]interface{}, 0)
	stopColumns := map[string]bool{
		"recruter_id": true,
		"id":          true,
		"name":        true,
		"ord":         true,
		"date_create": true,
	}

	isNeedAssert := oldData != nil && id == oldData.Id
	for _, col := range table.Columns() {
		name := col.Name()
		newVal := u.ColValue(name)
		if stopColumns[name] || EmptyValue(newVal) {
			continue
		}

		oldVal := oldData.ColValue(name)
		if (name == "user_ids") && (isNeedAssert || !reflect.DeepEqual(oldVal, newVal)) {
			columns = append(columns, name)
			args = append(args, newVal)
		} else if isNeedAssert || oldVal != newVal {
			columns = append(columns, name)
			args = append(args, newVal)
		}
	}

	i, err := table.Update(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(append(args, id)...),
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

		toLogVacancy(ctx, DB, u.SelectCompany.Id, id, text, CODE_LOG_UPDATE)
	}

	return createResult(i)
}

func HandleAddVacancy(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*VacancyDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	columns := []string{
		"platform_id",
		"seniority_id",
		"company_id",
		"location_id",
		"description",
		"details",
		"link",
		"status",
		"salary",
		"user_ids",
	}
	args := []interface{}{
		u.Platform_id,
		u.Seniority_id,
		u.Company_id,
		u.Location_id,
		u.Description,
		u.Details,
		u.Link,
		u.SelectedVacancyStatus,
		u.Salary,
		u.User_ids,
	}

	table, _ := db.NewVacancies(DB)
	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogVacancy(ctx, DB, u.SelectCompany.Id, int32(i), u.Description, CODE_LOG_INSERT)

	return createResult(i)
}

func HandleDeleteVacancy(ctx *fasthttp.RequestCtx) (interface{}, error) {
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

	table, _ := db.NewVacancies(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	err = DB.Conn.ExecDDL(ctx, "delete from vacancies where id = $1", id)
	if err != nil {
		return createErrResult(err)
	}

	toLogVacancy(ctx, DB, table.Record.Company_id, id, table.Record.Name.String, CODE_LOG_DELETE)
	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}

type DTOVacancy struct {
	CompanyId      int32 `json:"company_id"`
	WithRecruiters bool  `json:"withRecruiters"`
	IsActive       bool  `json:"isActive"`
}

func (d *DTOVacancy) GetValue() interface{} {
	return d
}

func (d *DTOVacancy) NewValue() interface{} {
	return &DTOVacancy{}
}

func toLogVacancy(ctx *fasthttp.RequestCtx, DB *dbEngine.DB, companyId, vacancyId int32, text string, code int32) {
	user := auth.GetUserData(ctx)
	toLog(ctx, DB,
		dbEngine.ColumnsForSelect("user_id", "company_id", "vacancy_id", "text", "date_create", "d_c",
			"kod_deystviya"),
		dbEngine.ArgsForSelect(user.Id, companyId, vacancyId,
			text,
			time.Now(),
			time.Now(),
			code))
}
