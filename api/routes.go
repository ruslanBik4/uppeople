// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"go/types"
	"runtime/trace"
	"time"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

var (
	ParamID = apis.InParam{
		Name:     "id",
		Desc:     "id of candidate",
		DefValue: apis.ApisValues(apis.ChildRoutePath),
		Req:      true,
		Type:     apis.NewTypeInParam(types.Int32),
	}
	ParamPageNum = apis.InParam{
		Name:     "page_num",
		Desc:     "number of page",
		DefValue: apis.ApisValues(apis.ChildRoutePath),
		Req:      true,
		Type:     apis.NewTypeInParam(types.Int),
	}
	ParamStartDate = apis.InParam{
		Name:     "start_date",
		Desc:     "first date of period",
		DefValue: apis.ApisValues(apis.ChildRoutePath),
		Req:      true,
		Type:     apis.NewTypeInParam(types.String),
	}
	ParamEndDate = apis.InParam{
		Name: "end_date",
		Desc: "last date of period",
		DefValue: func(ctx *fasthttp.RequestCtx) interface{} {
			return time.Now()
		},
		Req:  false,
		Type: apis.NewTypeInParam(types.String),
	}
	ParamCompanyID = apis.InParam{
		Name:     "company_id",
		Desc:     "id of company",
		DefValue: apis.ApisValues(apis.ChildRoutePath),
		Req:      true,
		Type:     apis.NewTypeInParam(types.Int32),
	}
)

var (
	Routes     = apis.NewMapRoutes()
	PostRoutes = apis.ApiRoutes{
		"/api/auth/login": {
			Fnc:  HandleAuthLogin,
			Desc: "login & return  candidate info",
			Params: []apis.InParam{
				{
					Name: "email",
					Type: apis.NewTypeInParam(types.String),
				},
				{
					Name: "password",
					Type: apis.NewTypeInParam(types.String),
				},
			},
			DTO:      &DTOAuth{},
			WithCors: true,
			// Resp:   search.RespGroups(),
		},
		"/api/main/addNewVacancy": {
			Fnc:      HandleAddVacancy,
			Desc:     "add new vacancy",
			DTO:      &VacancyDTO{},
			NeedAuth: true,
		},
		"/api/main/updateStatusCandidates": {
			Fnc:      HandleUpdateStatusCandidates,
			Desc:     "update status candidate",
			DTO:      &db.Vacancies_to_candidatesFields{},
			NeedAuth: true,
		},
		"/api/main/reContactCandidate/": {
			Fnc:      HandleReContactCandidate,
			Desc:     "update status candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
			WithCors: true,
		},
		"/api/main/addNewCompany": {
			Fnc:      HandleAddCompany,
			Desc:     "add new company",
			DTO:      &db.CompaniesFields{},
			NeedAuth: true,
		},
		"/api/main/editCompany/": {
			Fnc:      HandleEditCompany,
			Desc:     "edit new company",
			DTO:      &db.CompaniesFields{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/admin/editUser": {
			Fnc:      HandleEditUser,
			Desc:     "edit new users",
			DTO:      &DTOUser{},
			NeedAuth: true,
		},
		"/api/admin/add-platform": {
			Fnc:      HandleAddPlatform,
			Desc:     "add new platform",
			DTO:      &db.PlatformsFields{},
			NeedAuth: true,
		},
		"/api/main/addCommentForCompany/": {
			Fnc:      HandleAddCommentForCompany,
			Desc:     "add comment of company",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/addCommentForCandidate/": {
			Fnc:      HandleAddCommentsCandidate,
			Desc:     "add comment of candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/deleteCommentCandidate/": {
			Fnc:      HandleRmCommentsCandidate,
			Desc:     "delete comment of candidate (/id of comments)",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/interview/sendCV/": {
			Fnc:      HandleSendCV,
			Desc:     "send candidate to company",
			DTO:      &DTOSendCV{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/addNewContactForCompany/": {
			Fnc:      HandleAddContactForCompany,
			Desc:     "add contact of company",
			DTO:      &DTOContact{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/editContactCompany/": {
			Fnc:      HandleEditContactForCompany,
			Desc:     "edit contact of company",
			DTO:      &DTOContact{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/editVacancy/": {
			Fnc:      HandleEditVacancy,
			Desc:     "edit vacancy",
			DTO:      &VacancyDTO{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/updateStatusVacancy": {
			Fnc:      HandleEditStatusVacancy,
			Desc:     "edit vacancy",
			DTO:      &VacancyDTO{},
			NeedAuth: true,
		},
		"/api/main/addNewCandidate": {
			Fnc:      HandleAddCandidate,
			Desc:     "show search results according range of characteristics",
			DTO:      &CandidateDTO{},
			WithCors: true,
			NeedAuth: true,
		},
		"/api/main/addAvatarCandidate/": {
			Fnc:       HandleAddAvatar,
			Desc:      "put avatar img to photos",
			Multipart: true,
			NeedAuth:  true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/addLogoCompanies/": {
			Fnc:       HandleAddLogo,
			Desc:      "put logo to company",
			Multipart: true,
			NeedAuth:  true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/editCandidate/": {
			Fnc:      HandleEditCandidate,
			Desc:     "edit candidate",
			DTO:      &CandidateDTO{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/deleteVacancy/": {
			Fnc:      HandleDeleteVacancy,
			Desc:     "delete Vacancy",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/globalSearch": {
			Fnc:      HandleGlobalSearch,
			Desc:     "GlobalSearch",
			DTO:      &GlobalSearch{},
			NeedAuth: true,
		},
		// todo profile vacancy
		"/api/admin/view-editUser/": {
			Fnc:      HandleGetUser,
			Desc:     "get Users data",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/admin/newUser/": {
			Fnc:      HandleNewUser,
			Desc:     "create new user",
			DTO:      &DTOUser{},
			NeedAuth: true,
		},
		"/api/admin/deleteUser/": {
			Fnc:      HandleDelUser,
			Desc:     "delete User",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/deleteCandidate/": {
			Fnc:      HandleDeleteCandidate,
			Desc:     "delete candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/saveFollowUp": {
			Fnc:      HandleFollowUpCandidate,
			Desc:     "show search results according range of characteristics",
			DTO:      &FollowUpDTO{},
			NeedAuth: true,
		},
		"/api/main/returnAllVacancy": {
			Fnc:      HandleReturnAllVacancy,
			Desc:     "get list of vacancies",
			DTO:      &DTOVacancy{},
			NeedAuth: true,
		},
		"/api/main/viewAllVacancyInCompany/": {
			Fnc:      HandleViewAllVacancyInCompany,
			Desc:     "get list of vacancies",
			DTO:      &vacDTO{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/": {
			Fnc:  HandleApiRedirect,
			Desc: "show search results according range of characteristics",
			// DTO:    &DTOSearch{},
			Method:   apis.POST,
			NeedAuth: true,
			// Resp:   search.RespGroups(),
		},
		"/api/main/getTags": {
			Fnc:      HandleGetTags,
			Desc:     "return select of tags",
			Method:   apis.POST,
			NeedAuth: true,
		},
		"/api/main/getStatuses": {
			Fnc:      HandleGetStatuses,
			Desc:     "return select of tags",
			Method:   apis.POST,
			NeedAuth: true,
		},
		"/api/main/allCandidates/": {
			Fnc:      HandleAllCandidates,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &SearchCandidates{},
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/returnAllCandidates/": {
			Fnc:      HandleReturnAllCandidate,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &SearchCandidates{},
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/viewCandidatesFreelancerOnVacancies/": {
			Fnc:      HandleCandidatesFreelancerOnVacancies,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &SearchCandidates{},
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/viewAllCandidatesForCompany/": {
			Fnc:      HandleAllCandidatesForCompany,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &SearchCandidates{},
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/returnAllCompanies/": {
			Fnc:      HandleAllCompanies,
			Desc:     "show all companies",
			NeedAuth: true,
			DTO:      &SearchCompany{},
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/getCandidatesAmountByTags": {
			Fnc:      HandleGetCandidatesAmountByTags,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &DTOAmounts{},
		},
		"/api/reports/by_tags": {
			Fnc:      HandleDownloadReportByTag,
			Desc:     "reports by tags",
			NeedAuth: true,
			DTO:      &DTOAmounts{},
		},
		"/api/reports/by_status": {
			Fnc:      HandleDownloadReportByStatus,
			Desc:     "reports by tags",
			NeedAuth: true,
			DTO:      &DTOAmounts{},
		},
		"/api/main/getCandidatesAmountByStatuses": {
			Fnc:      HandleGetCandidatesAmountByStatuses,
			Desc:     "show all candidates",
			NeedAuth: true,
			DTO:      &DTOAmounts{},
		},
		"/api/main/getCandidatesAmountByVacancies": {
			Fnc:      HandleGetCandidatesByVacancies,
			Desc:     "show all GetCandidatesByVacancies",
			NeedAuth: true,
			DTO:      &DTOAmounts{},
		},
		"/api/interview/inviteOnInterviewSend/": {
			Fnc:      HandleInviteOnInterviewSend,
			Desc:     "invite to interview candidate",
			DTO:      &DTOSendInterview{},
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
	}
	GetRoutes = apis.ApiRoutes{
		"/api/img/": {
			Fnc:  HandleGetImg,
			Desc: "show img of company",
		},
		"/img/photos/": {
			Fnc:  HandleImg,
			Desc: "show img of candidate",
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/reports/": {
			Fnc:  HandleDownloadWholeReport,
			Desc: "reports by tags of all companies, vacancies, recruiters",
			Params: []apis.InParam{
				ParamStartDate,
				ParamEndDate,
			},
		},
		"/api/main/dashBoardAdmin": {
			Fnc:      HandleDashBoard,
			Desc:     "show HandleDashBoard for admin",
			NeedAuth: true,
		},
		"/api/main/dashBoardRecruiter": {
			Fnc:      HandleDashBoard,
			Desc:     "show HandleDashBoard for recruiter",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/viewInformationForCompany/": {
			Fnc:      HandleInformationForCompany,
			Desc:     "show company by $id",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/viewEditContactCompany/": {
			Fnc:      HandleViewContactForCompany,
			Desc:     "show contacts of company  by $id_contacts",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/deleteContactForCompany/": {
			Fnc:      HandleDeleteContactForCompany,
			Desc:     "remove contacts of company  by $id_contacts",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/commentsCompany/": {
			Fnc:      HandleCommentsCompany,
			Desc:     "show company by $id",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/returnAllCandidates/": {
			Fnc:      HandleReturnAllCandidate,
			Desc:     "show all candidates",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/viewCandidatesFreelancerOnVacancies/": {
			Fnc: func(ctx *fasthttp.RequestCtx) (interface{}, error) {
				ctx.SetStatusCode(fasthttp.StatusNoContent)
				return nil, nil
			},
			Desc:     "show all candidates",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/admin/all-staff": {
			Fnc:      HandleAllStaff,
			Desc:     "show all users",
			NeedAuth: true,
		},
		"/api/admin/returnAllPlatforms/": {
			Fnc:      HandleAllPlatforms,
			Desc:     "show all users",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/returnOptionsForSelects": {
			Fnc:      HandleReturnOptionsForSelects,
			Desc:     "show options for selects ",
			NeedAuth: true,
			WithCors: true,
		},
		"/api/main/allCandidates/": {
			Fnc:      HandleAllCandidates,
			Desc:     "show all candidates",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamPageNum,
			},
		},
		"/api/main/viewOneCandidate/": {
			Fnc:      HandleViewCandidate,
			Desc:     "show one candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/interview/viewInformationForSendCV/": {
			Fnc:      HandleInformationForSendCV,
			Desc:     "get data to sendCV candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/interview/inviteOnInterviewView/": {
			Fnc:      HandleInviteOnInterviewView,
			Desc:     "invite to interview candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/commentsCandidate/": {
			Fnc:      HandleCommentsCandidate,
			Desc:     "show comments for candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/main/viewVacancy/": {
			Fnc:      HandleViewVacancy,
			Desc:     "show one candidate",
			NeedAuth: true,
			Params: []apis.InParam{
				ParamID,
			},
		},
		"/api/admin/returnLogsForCand/": apis.NewAPIRouteWithDBEngine(
			"show logs of candidate",
			apis.GET,
			true,
			[]apis.InParam{
				ParamID,
				{
					Name:     "isCand",
					DefValue: true,
					Type:     apis.NewTypeInParam(types.Bool),
				},
			},
			"get_log"),
		"/api/admin/returnLogsForCompany/": apis.NewAPIRouteWithDBEngine(
			"show logs of company",
			apis.GET,
			true,
			[]apis.InParam{
				ParamCompanyID,
				{
					Name:     "isCand",
					DefValue: false,
					Type:     apis.NewTypeInParam(types.Bool),
				},
			},
			"get_log"),
		// "/api/main/returnAllCandidates/": {
		// 	Fnc:      HandleAllCandidates,
		// 	Desc:     "show search results according range of characteristics",
		// 	NeedAuth: true,
		// },
		// "/api/main/viewCandidatesFreelancerOnVacancies/": {
		// 	Fnc:      HandleAllCandidates,
		// 	Desc:     "show search results according range of characteristics",
		// 	NeedAuth: true,
		// },
		"/api/": {
			Fnc:      HandleApiRedirect,
			Desc:     "show search results according range of characteristics",
			NeedAuth: true,
			// DTO:    &DTOSearch{},
			// Method: apis.POST,
			// Resp:   search.RespGroups(),
		},
		"/api/system/stat_conn/": {
			Fnc: HandleStatConn,
			Params: []apis.InParam{
				{
					Name: "cmd",
					Desc: "command for pgx",
					Req:  false,
					Type: apis.NewTypeInParam(types.String),
				},
			},
		},
		"/api/system/trace/": {
			Fnc: HandleTrace,
			Params: []apis.InParam{
				{
					Name: "cmd",
					Desc: "command for pgx",
					Req:  false,
					Type: apis.NewTypeInParam(types.String),
				},
			},
		},
	}
)
var hosts = []string{
	"back.uppeople.co",
	"back2.uppeople.co",
}

func doRequest(req *fasthttp.Request, resp *fasthttp.Response, host string) error {
	c := fasthttp.Client{}
	uri := req.URI()
	uri.SetScheme("http")
	uri.SetHost(host)

	err := c.Do(req, resp)
	if err != nil {
		return errors.Wrap(err, "do")
	}

	return nil
}

func HandleApiRedirect(ctx *fasthttp.RequestCtx) (interface{}, error) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	return nil, nil
}

func init() {
	for path, route := range PostRoutes {
		route.Method = apis.POST
		if trace.IsEnabled() {
			route.Fnc = func(handler apis.ApiRouteHandler) apis.ApiRouteHandler {
				return func(ctx *fasthttp.RequestCtx) (resp interface{}, err error) {
					ctx1, task := trace.NewTask(ctx, path)
					defer task.End()
					reg := trace.StartRegion(ctx1, path)
					defer reg.End()
					logs.DebugLog(reg, task)
					trace.WithRegion(ctx1, path,
						func() {
							resp, err = handler(ctx)
						})

					return
				}
			}(route.Fnc)
		}
	}
	Routes.AddRoutes(PostRoutes)

	for path, route := range GetRoutes {
		if trace.IsEnabled() {
			route.Fnc = func(handler apis.ApiRouteHandler) apis.ApiRouteHandler {
				return func(ctx *fasthttp.RequestCtx) (resp interface{}, err error) {
					ctx1, task := trace.NewTask(ctx, path)
					defer task.End()
					reg := trace.StartRegion(ctx1, path)
					defer reg.End()
					logs.DebugLog(reg, task)
					trace.WithRegion(ctx1, path,
						func() {
							resp, err = handler(ctx)
						})

					return
				}
			}(route.Fnc)
		}
	}

	Routes.AddRoutes(GetRoutes)
}
