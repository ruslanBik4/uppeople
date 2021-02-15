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
	CompanyID             int32          `json:"company_id"`
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

type VacanciesView struct {
	*db.VacanciesFields
	Date      string `json:"date"`
	Platform  string `json:"platform"`
	Company   string `json:"company"`
	Location  string `json:"location"`
	Seniority string `json:"seniority"`
}

func (v *VacanciesView) GetFields(columns []dbEngine.Column) []interface{} {
	res := make([]interface{}, len(columns))
	for i, col := range columns {
		switch col.Name() {
		case "platform":
			res[i] = v.Platform
		case "company":
			res[i] = v.Company
		case "location":
			res[i] = v.Location
		case "seniority":
			res[i] = v.Seniority
		default:
			res[i] = v.RefColValue(col.Name())
		}
	}

	return res
}

type ResVacancies struct {
	*ResList
	CandidateStatus SelectedUnits   `json:"candidateStatus"`
	VacancyStatus   SelectedUnits   `json:"vacancyStatus"`
	Vacancies       []VacanciesView `json:"vacancies"`
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
			(select s.nazva from seniority s where v.seniority_id=s.id) as seniority,
			(select c.name from company c where v.company_id=c.id) as company,
			(select s.name from location_for_vacancies s where v.location_id=s.id) as location
			from vacancies v
			where id = $1
`,
		ctx.UserValue("id"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "	")
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
		toLogVacancy(ctx, DB, u.SelectCompany.Id, u.Id, "status", 100)
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

	u.User_ids = ""
	for i, unit := range u.SelectRecruiter {
		if i > 0 {
			u.User_ids += ", "
		}
		u.User_ids += fmt.Sprintf("%d", unit.Id)
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
	}
	u.Platform_id = u.SelectPlatform.Id
	u.Seniority_id = u.SelectSeniority.Id
	u.Company_id = u.SelectCompany.Id
	u.Location_id = u.SelectLocation.Id
	u.Status = u.SelectedVacancyStatus
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
		args = []interface{}{
			u.SelectPlatform.Id,
			u.SelectSeniority.Id,
			u.SelectCompany.Id,
			u.SelectLocation.Id,
			u.Description,
			u.Details,
			u.Link,
			u.SelectedVacancyStatus,
			u.Salary,
			u.User_ids,
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

	userIDs, comma := "", ""
	for _, unit := range u.SelectRecruiter {
		userIDs += fmt.Sprintf("%s%d", comma, unit.Id)
		comma = "-"
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
		u.SelectPlatform.Id,
		u.SelectSeniority.Id,
		u.SelectCompany.Id,
		u.SelectLocation.Id,
		u.Description,
		u.Details,
		u.Link,
		u.SelectedVacancyStatus,
		u.Salary,
		userIDs,
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

	toLogVacancy(ctx, DB, table.Record.Company_id, id, table.Record.Name.String, 103)
	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}

func HandleViewAllVacancyInCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}
	offset := 0
	id, ok := ctx.UserValue(ParamPageNum.Name).(int)
	if ok && id > 1 {
		offset = id * pageItem
	}

	filter, ok := ctx.UserValue(apis.JSONParams).(*vacDTO)
	if !ok {
		return "DTO is wrong", apis.ErrWrongParamsList
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)
	if filter.CompanyID > 0 {
		columns = append(columns, "company_id")
		args = append(args, filter.CompanyID)
	}

	if l := len(filter.SelectStatuses); l > 0 {
		arg := make([]int32, l)
		for i, s := range filter.SelectStatuses {
			arg[i] = s.Id
		}
		columns = append(columns, "status")
		args = append(args, arg)
	}

	if l := len(filter.SelectPlatforms); l > 0 {
		arg := make([]int32, l)
		for i, s := range filter.SelectPlatforms {
			arg[i] = s.Id
		}
		columns = append(columns, "platform_id")
		args = append(args, arg)
	}

	// if l := len(filter.SelectCandidateStatus); l > 0 {
	// 	arg := make([]int32, l)
	// 	for i, s := range filter.SelectCandidateStatus {
	// 		arg[i] = s.Id
	// 	}
	// 	columns = append(columns, "status")
	// 	args = append(args, arg)
	// }
	//
	if l := len(filter.SelectSeniorities); l > 0 {
		arg := make([]int32, l)
		for i, s := range filter.SelectSeniorities {
			arg[i] = s.Id
		}
		columns = append(columns, "seniority_id")
		args = append(args, arg)
	}

	vacancies, _ := db.NewVacancies(DB)
	res := ResVacancies{
		ResList:         NewResList(ctx, DB, id),
		Vacancies:       make([]VacanciesView, 0),
		CandidateStatus: getStatusVac(ctx, DB),
		VacancyStatus:   getStatuses(ctx, DB),
	}

	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy("date_create desc"),
		dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	}
	if len(columns) > 0 {
		options = append(options, dbEngine.WhereForSelect(columns...))
		options = append(options, dbEngine.ArgsForSelect(args...))
	}
	i := 0
	companies, _ := db.NewCompanies(DB)
	locs, _ := db.NewLocation_for_vacancies(DB)
	err := vacancies.SelectSelfScanEach(ctx,
		func(record *db.VacanciesFields) error {
			view := VacanciesView{
				VacanciesFields: record,
				Date:            record.Date_create.Format("2006-01-02"),
				Company:         "",
				Location:        "",
			}
			if record.Company_id > 0 {
				err := companies.SelectOneAndScan(ctx,
					&view.Company,
					dbEngine.ColumnsForSelect("name"),
					dbEngine.WhereForSelect("id"),
					dbEngine.ArgsForSelect(record.Company_id),
				)
				if err != nil {
					logs.ErrorLog(err, "companies.SelectOneAndScan")
				}
			}

			if record.Location_id > 0 {
				err := locs.SelectOneAndScan(ctx,
					&view.Location,
					dbEngine.ColumnsForSelect("name"),
					dbEngine.WhereForSelect("id"),
					dbEngine.ArgsForSelect(record.Location_id),
				)
				if err != nil {
					logs.ErrorLog(err, "locs.SelectOneAndScan")
				}

			}

			for _, s := range res.Seniority {
				if s.Id == int32(record.Seniority_id) {
					view.Seniority = s.Label
					break
				}
			}

			for _, s := range res.Platforms {
				if s.Id == record.Platform_id {
					view.Platform = s.Label
					break
				}
			}

			res.Vacancies = append(res.Vacancies, view)

			i++
			if i == pageItem {
				return errLimit
			}

			return nil
		},
		options...,
	)

	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	return res, nil

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
