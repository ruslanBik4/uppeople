// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

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
			res[i] = &v.Platform
		case "company":
			res[i] = &v.Company
		case "location":
			res[i] = &v.Location
		case "seniority":
			res[i] = &v.Seniority
		default:
			res[i] = v.RefColValue(col.Name())
		}
	}

	return res
}

type ResVacancies struct {
	*ResList
	CandidateStatus db.SelectedUnits `json:"candidateStatus"`
	VacancyStatus   db.SelectedUnits `json:"vacancyStatus"`
	Vacancies       []VacanciesView  `json:"vacancies"`
}

func HandleReturnAllVacancy(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	dto, ok := ctx.UserValue(apis.JSONParams).(*DTOVacancy)
	if !ok {
		return "DTO is wrong", apis.ErrWrongParamsList
	}

	where, comma := "", "where"
	if p := dto.CompanyId; p > 0 {
		where += fmt.Sprintf(" %s v.company_id = %d", comma, p)
		comma = "AND"
	}

	if dto.IsActive {
		where += comma + " v.status = ANY(array[0, 1])"
	}

	sql := `select v.id, company_id, CONCAT(c.name, ' (', platforms.nazva, ') ',
		(select s.status from statuses s where s.id = v.status)
) as name`
	if dto.WithRecruiters {
		sql += `, (SELECT array_agg(distinct user_id) as recruiter_id
                            FROM vacancies_to_candidates
                            WHERE vacancy_id = v.id) as recruiters`
	}
	sql += ` from vacancies v left join platforms on v.platform_id=platforms.id
	left join companies c on v.company_id = c.id
`
	return DB.Conn.SelectToMaps(ctx, sql+where+" order by v.status")
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

	dto, ok := ctx.UserValue(apis.JSONParams).(*vacDTO)
	if !ok {
		return "DTO is wrong", apis.ErrWrongParamsList
	}

	columns := make([]string, 0)
	args := make([]interface{}, 0)
	if dto.CompanyId > 0 {
		columns = append(columns, "company_id")
		args = append(args, dto.CompanyId)
	}

	if l := len(dto.SelectStatuses); l > 0 {
		arg := make([]int32, l)
		for i, s := range dto.SelectStatuses {
			arg[i] = s.Id
		}
		columns = append(columns, "status")
		args = append(args, arg)
	}

	if l := len(dto.SelectPlatforms); l > 0 {
		arg := make([]int32, l)
		for i, s := range dto.SelectPlatforms {
			arg[i] = s.Id
		}
		columns = append(columns, "platform_id")
		args = append(args, arg)
	}

	// if l := len(dto.SelectCandidateStatus); l > 0 {
	// 	arg := make([]int32, l)
	// 	for i, s := range dto.SelectCandidateStatus {
	// 		arg[i] = s.Id
	// 	}
	// 	columns = append(columns, "status")
	// 	args = append(args, arg)
	// }
	//
	if l := len(dto.SelectSeniorities); l > 0 {
		arg := make([]int32, l)
		for i, s := range dto.SelectSeniorities {
			arg[i] = s.Id
		}
		columns = append(columns, "seniority_id")
		args = append(args, arg)
	}

	orderBy := "date_create desc"
	if dto.CurrentColumn > "" {

		switch dto.CurrentColumn {
		case "Company":
			orderBy = `(select name from companies where id = company_id)`
		case "Platform":
			orderBy = `(select nazva from platforms where id = platform_id)`
		case "Location":
			orderBy = `(select name from location_for_vacancies where id = location_id)`
		case "Seniority":
			orderBy = `(select nazva from seniorities where id = seniority_id)`
		case "Contacts":
			orderBy = `coalesce(email, phone, skype)`
		case "Date":
			// orderBy = `coalesce(email, phone, skype)`
		default:
			orderBy = dto.CurrentColumn
		}
		if dto.Sort > 0 {
			orderBy += " desc"
		}
	}

	vacancies, _ := db.NewVacancies(DB)
	res := ResVacancies{
		ResList:         NewResList(id),
		Vacancies:       make([]VacanciesView, 0),
		CandidateStatus: db.GetStatusForVacAsSelectedUnits(),
		VacancyStatus:   db.GetStatusAsSelectedUnits(),
	}

	options := []dbEngine.BuildSqlOptions{
		dbEngine.OrderBy(orderBy),
		dbEngine.FetchOnlyRows(pageItem),
		dbEngine.Offset(offset),
	}
	optionsCount := []dbEngine.BuildSqlOptions{
		dbEngine.ColumnsForSelect("count(*)"),
	}

	if len(columns) > 0 {
		options = append(options, dbEngine.WhereForSelect(columns...), dbEngine.ArgsForSelect(args...))
		optionsCount = append(optionsCount, dbEngine.WhereForSelect(columns...), dbEngine.ArgsForSelect(args...))
	}

	companies, _ := db.NewCompanies(DB)
	locs, _ := db.NewLocation_for_vacancies(DB)
	err := vacancies.SelectSelfScanEach(ctx,
		func(record *db.VacanciesFields) error {
			view := VacanciesView{
				VacanciesFields: record,
				Date:            record.DateCreate.Format("2006-01-02"),
				Company:         "",
				Location:        "",
			}
			if record.CompanyId > 0 {
				err := companies.SelectOneAndScan(ctx,
					&view.Company,
					dbEngine.ColumnsForSelect("name"),
					dbEngine.WhereForSelect("id"),
					dbEngine.ArgsForSelect(record.CompanyId),
				)
				if err != nil {
					logs.ErrorLog(err, "companies.SelectOneAndScan")
				}
			}

			if record.LocationId > 0 {
				err := locs.SelectOneAndScan(ctx,
					&view.Location,
					dbEngine.ColumnsForSelect("name"),
					dbEngine.WhereForSelect("id"),
					dbEngine.ArgsForSelect(record.LocationId),
				)
				if err != nil {
					logs.ErrorLog(err, "locs.SelectOneAndScan")
				}

			}

			view.Seniority = db.GetSeniorityFromId(record.SeniorityId).Nazva.String

			view.Platform = db.GetPlatformFromId(record.PlatformId).Nazva.String

			res.Vacancies = append(res.Vacancies, view)

			return nil
		},
		options...,
	)

	if err != nil {
		return nil, errors.Wrap(err, "	")
	}

	if len(res.Vacancies) < pageItem {
		res.ResList.TotalPage = 1
		res.ResList.Count = len(res.Vacancies)
	} else {
		err = vacancies.SelectOneAndScan(ctx,
			&res.ResList.Count,
			optionsCount...)
		if err != nil {
			logs.ErrorLog(err, "count")
		} else {
			res.ResList.TotalPage = res.ResList.Count / pageItem
		}
	}

	return res, nil
}
