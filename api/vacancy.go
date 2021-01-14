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

	table, _ := db.NewVacancies(DB)
	err := table.SelectOneAndScan(ctx,
		table,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(ctx.UserValue("id")),
	)
	if err != nil && err != errLimit {
		return nil, errors.Wrap(err, "	")
	}

	return table.Record, nil
}

func HandleAddVacancy(ctx *fasthttp.RequestCtx) (interface{}, error) {
	u, ok := ctx.UserValue(apis.JSONParams).(*VacancyDTO)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	logs.DebugLog(u)
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
	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if ok {
		i, err := table.Update(ctx,
			dbEngine.ColumnsForSelect(columns...),
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(append(args, id)...),
		)
		if err != nil {
			return createErrResult(err)
		}

		if i > 0 {
			toLogVacancy(ctx, DB, u.SelectCompany.Id, id, u.Description, 100)
		}

		return createResult(i)
	}

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
		dbEngine.WhereForSelect("company_id", "platform_id", "seniority_id"),
		dbEngine.ArgsForSelect(u.SelectCompany.Id, u.SelectPlatform.Id, u.SelectSeniority.Id),
	)
	if err != nil {
		logs.ErrorLog(err, "table.SelectOneAndScan")
	}

	toLogVacancy(ctx, DB, u.SelectCompany.Id, u.Id, u.Description, 101)

	return createResult(i)
}

func HandleViewAllVacancyInCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	filter, ok := ctx.UserValue(apis.JSONParams).(*vacDTO)
	if !ok {
		return "DTO is wrong", apis.ErrWrongParamsList
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)
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
		ResList:         NewResList(ctx, DB),
		Vacancies:       make([]VacanciesView, 0),
		CandidateStatus: getStatusVac(ctx, DB),
		VacancyStatus:   getStatuses(ctx, DB),
	}

	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy("date_create desc"),
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
				Date:            record.Date_create.Format("2002-01-02"),
				Company:         "",
				Location:        "",
			}
			if record.Company_id.Valid {
				err := companies.SelectOneAndScan(ctx,
					&view.Company,
					dbEngine.ColumnsForSelect("name"),
					dbEngine.WhereForSelect("id"),
					dbEngine.ArgsForSelect(record.Company_id.Int64),
				)
				if err != nil {
					logs.ErrorLog(err, "companies.SelectOneAndScan")
				}
			}

			if record.Location_id.Valid {
				err := locs.SelectOneAndScan(ctx,
					&view.Location,
					dbEngine.ColumnsForSelect("name"),
					dbEngine.WhereForSelect("id"),
					dbEngine.ArgsForSelect(record.Location_id.Int64),
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
				if s.Id == int32(record.Platform_id.Int64) {
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
