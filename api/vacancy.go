// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v4"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

type VacancyDTO struct {
	*db.VacanciesFields
	Comment               string            `json:"comment"`
	SelectCompany         db.SelectedUnit   `json:"selectCompany"`
	SelectLocation        db.SelectedUnit   `json:"selectLocation"`
	SelectPlatform        db.SelectedUnit   `json:"selectPlatform"`
	SelectSeniority       db.SelectedUnit   `json:"selectSeniority"`
	SelectRecruiter       []db.SelectedUnit `json:"selectRecruiter"`
	SelectedVacancyStatus int32             `json:"selectedVacancyStatus"`
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
	CompanyId             int32             `json:"company_id"`
	Sort                  int32             `json:"sort"`
	CurrentColumn         string            `json:"currentColumn"`
	SelectPlatforms       []db.SelectedUnit `json:"selectPlatforms"`
	SelectSeniorities     []db.SelectedUnit `json:"selectSeniorities"`
	SelectCandidateStatus []db.SelectedUnit `json:"selectCandidate_status"`
	SelectStatuses        []db.SelectedUnit `json:"selectStatuses"`
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
			(select p.name from public.platforms p where v.platform_id=p.id) as platform,
			(select s.name from public.seniorities s where v.seniority_id=s.id) as seniority,
			(select c.name from companies c where v.company_id=c.id) as company,
			(select s.name from public.location_for_vacancies s where v.location_id=s.id) as location
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

	v.Date = v.DateCreate.Format("2006-01-02")

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
		toLogVacancyUpdate(ctx, u.SelectCompany.Id, u.Id, map[string]interface{}{"status": u.Id})
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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
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

		if isNeedAssert {
			oldVal := oldData.ColValue(name)
			if (name == "user_ids" && reflect.DeepEqual(oldVal, newVal)) ||
				(name != "user_ids" && oldVal == newVal) {
				continue
			}
		}

		columns = append(columns, name)
		args = append(args, newVal)
	}
	if len(columns) == 0 {
		return "no new data on record", apis.ErrWrongParamsList
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
		toLogVacancyUpdate(ctx, u.SelectCompany.Id, id, toLogUpdateValues(columns, args))
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
		u.PlatformId,
		u.SeniorityId,
		u.CompanyId,
		u.LocationId,
		u.Description,
		u.Details,
		u.Link,
		u.SelectedVacancyStatus,
		u.Salary,
		u.UserIds,
	}

	table, _ := db.NewVacancies(DB)
	i, err := table.Insert(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	toLogVacancyInsert(ctx, u.SelectCompany.Id, int32(i), u.Description)

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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
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

	toLogVacancyDelete(ctx, table.Record.CompanyId, id, table.Record.Name.String)
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
