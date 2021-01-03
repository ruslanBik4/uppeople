// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
)

type DTOAuth struct {
	Email, Password string
}

func (a *DTOAuth) GetValue() interface{} {
	return a
}

func (a *DTOAuth) NewValue() interface{} {
	return &DTOAuth{}
}

func HandleAuthLogin(ctx *fasthttp.RequestCtx) (interface{}, error) {
	a, ok := ctx.UserValue(apis.JSONParams).(*DTOAuth)
	if ok && a.Email == "test@test.com" && a.Password == "1111" {
		v := map[string]interface{}{
			"access_token": "EElR0VdJnJzianbkx5y4JQ==",
			"expires_at":   "2022-01-01 19:43:20",
			"optionsForSelects": map[string]interface{}{
				"candidateStatus": []map[string]interface{}{
					map[string]interface{}{
						"id":    1,
						"label": "Before pre-screen",
						"value": "before pre-screen",
					},
					map[string]interface{}{
						"id":    4,
						"label": "Final Interview",
						"value": "final interview",
					},
					map[string]interface{}{
						"id":    6,
						"label": "Hired",
						"value": "hired",
					},
					map[string]interface{}{
						"id":    2,
						"label": "Interview",
						"value": "interview",
					},
					map[string]interface{}{
						"id":    5,
						"label": "OFFER",
						"value": "offer",
					},
					map[string]interface{}{
						"id":    11,
						"label": "On hold",
						"value": "on hold",
					},
					map[string]interface{}{
						"id":    13,
						"label": "Pre-screen",
						"value": "pre-screen",
					},
					map[string]interface{}{
						"id":    10,
						"label": "Rejected",
						"value": "rejected",
					},
					map[string]interface{}{
						"id":    9,
						"label": "Review",
						"value": "review",
					},
					map[string]interface{}{
						"id":    3,
						"label": "Test",
						"value": "test",
					},
					map[string]interface{}{
						"id":    8,
						"label": "WR",
						"value": "wr",
					},
				},
				"companies": []map[string]interface{}{
					{
						"id":    2,
						"label": "Autoreturn (Zaly)",
						"value": "autoreturn (zaly)",
					},
					{
						"id":    29,
						"label": "BSG",
						"value": "bsg",
					},
					{
						"id":    1,
						"label": "Cars West",
						"value": "cars west",
					},
					{
						"id":    31,
						"label": "CMK",
						"value": "cmk",
					},
					{
						"id":    3,
						"label": "Dev-Pro.net",
						"value": "dev-pro.net",
					},
					{
						"id":    4,
						"label": "Diatom Enterprises",
						"value": "diatom enterprises",
					},
					{
						"id":    5,
						"label": "DiJust Development",
						"value": "dijust development",
					},
					{
						"id":    14,
						"label": "Edible Arrangements",
						"value": "edible arrangements",
					},
					{
						"id":    6,
						"label": "EnvisionTEC GmbH",
						"value": "envisiontec gmbh",
					},
					{
						"id":    30,
						"label": "FINEKO",
						"value": "fineko",
					},
					{
						"id":    28,
						"label": "FitGrid",
						"value": "fitgrid",
					},
					{
						"id":    7,
						"label": "GIS Art",
						"value": "gis art",
					},
					{
						"id":    8,
						"label": "iLogos Game Studios",
						"value": "ilogos game studios",
					},
					{
						"id":    26,
						"label": "inCust",
						"value": "incust",
					},
					{
						"id":    9,
						"label": "INSART",
						"value": "insart",
					},
					{
						"id":    10,
						"label": "iX.co",
						"value": "ix.co",
					},
					{
						"id":    33,
						"label": "Kuna",
						"value": "kuna",
					},
					{
						"id":    11,
						"label": "LALAFO",
						"value": "lalafo",
					},
					{
						"id":    36,
						"label": "moneyveo",
						"value": "moneyveo",
					},
					{
						"id":    12,
						"label": "MWDN ltd.",
						"value": "mwdn ltd.",
					},
					{
						"id":    34,
						"label": "MyCredit",
						"value": "mycredit",
					},
					{
						"id":    35,
						"label": "MyCredit",
						"value": "mycredit",
					},
					{
						"id":    13,
						"label": "MyHeritage",
						"value": "myheritage",
					},
					{
						"id":    15,
						"label": "Newxel",
						"value": "newxel",
					},
					{
						"id":    16,
						"label": "Noam Labs",
						"value": "noam labs",
					},
					{
						"id":    17,
						"label": "Pragmatic Play",
						"value": "pragmatic play",
					},
					{
						"id":    18,
						"label": "RevJet",
						"value": "revjet",
					},
					{
						"id":    19,
						"label": "RoboPlan",
						"value": "roboplan",
					},
					{
						"id":    20,
						"label": "SAMAWATT SA",
						"value": "samawatt sa",
					},
					{
						"id":    21,
						"label": "SkyDigitalLab",
						"value": "skydigitallab",
					},
					{
						"id":    27,
						"label": "SPD-Ukraine",
						"value": "spd-ukraine",
					},
					{
						"id":    22,
						"label": "Totango",
						"value": "totango",
					},
					{
						"id":    23,
						"label": "Trust Sourcing",
						"value": "trust sourcing",
					},
					{
						"id":    32,
						"label": "Verbit.ai",
						"value": "verbit.ai",
					},
					{
						"id":    24,
						"label": "Visonic (Tyco International)",
						"value": "visonic (tyco international)",
					},
					{
						"id":    37,
						"label": "WeAR Studio",
						"value": "wear studio",
					},
					{
						"id":    25,
						"label": "Zoolatech",
						"value": "zoolatech",
					},
				},
				"location": []map[string]interface{}{
					{
						"id":    6,
						"label": "Dnipro",
						"value": "dnipro",
					},
					{
						"id":    5,
						"label": "Kharkiv",
						"value": "kharkiv",
					},
					{
						"id":    2,
						"label": "Kyiv",
						"value": "kyiv",
					},
					{
						"id":    3,
						"label": "Lviv",
						"value": "lviv",
					},
					{
						"id":    4,
						"label": "Odessa",
						"value": "odessa",
					},
					{
						"id":    1,
						"label": "Remote",
						"value": "remote",
					},
				},
				"platforms": []map[string]interface{}{
					{
						"id":    5,
						"label": ".Net",
						"value": ".net",
					},
					{
						"id":    101,
						"label": ".NET+Angular",
						"value": ".net+angular",
					},
					{
						"id":    12,
						"label": "2D Animator",
						"value": "2d animator",
					},
					{
						"id":    50,
						"label": "2D Artist",
						"value": "2d artist",
					},
					{
						"id":    23,
						"label": "3D Artist",
						"value": "3d artist",
					},
					{
						"id":    16,
						"label": "Admin",
						"value": "admin",
					},
					{
						"id":    15,
						"label": "Android",
						"value": "android",
					},
					{
						"id":    36,
						"label": "Backend",
						"value": "backend",
					},
					{
						"id":    99,
						"label": "BIGDATA ENGINEER",
						"value": "bigdata engineer",
					},
					{
						"id":    53,
						"label": "Blockchain/Crypto currency devel",
						"value": "blockchain/crypto currency devel",
					},
					{
						"id":    30,
						"label": "Business Analyst",
						"value": "business analyst",
					},
					{
						"id":    88,
						"label": "Business Analyst (Post Launch)",
						"value": "business analyst (post launch)",
					},
					{
						"id":    89,
						"label": "C#",
						"value": "c#",
					},
					{
						"id":    24,
						"label": "C++",
						"value": "c++",
					},
					{
						"id":    38,
						"label": "Computational Linguist",
						"value": "computational linguist",
					},
					{
						"id":    80,
						"label": "Concept Artist",
						"value": "concept artist",
					},
					{
						"id":    21,
						"label": "CTO",
						"value": "cto",
					},
					{
						"id":    97,
						"label": "Data Analyst",
						"value": "data analyst",
					},
					{
						"id":    64,
						"label": "Data Architect",
						"value": "data architect",
					},
					{
						"id":    83,
						"label": "Data Engineer",
						"value": "data engineer",
					},
					{
						"id":    29,
						"label": "Data Scientist",
						"value": "data scientist",
					},
					{
						"id":    87,
						"label": "Database Administrator",
						"value": "database administrator",
					},
					{
						"id":    91,
						"label": "DB developer",
						"value": "db developer",
					},
					{
						"id":    63,
						"label": "Delivery Manager",
						"value": "delivery manager",
					},
					{
						"id":    81,
						"label": "Delphi",
						"value": "delphi",
					},
					{
						"id":    100,
						"label": "Deployment Engineer",
						"value": "deployment engineer",
					},
					{
						"id":    17,
						"label": "DevOps",
						"value": "devops",
					},
					{
						"id":    52,
						"label": "Director of Software Engineering",
						"value": "director of software engineering",
					},
					{
						"id":    71,
						"label": "Embedded C",
						"value": "embedded c",
					},
					{
						"id":    58,
						"label": "Engineering Manager",
						"value": "engineering manager",
					},
					{
						"id":    98,
						"label": "Flutter",
						"value": "flutter",
					},
					{
						"id":    2,
						"label": "Front-end",
						"value": "front-end",
					},
					{
						"id":    78,
						"label": "Front-End / Angular",
						"value": "front-end / angular",
					},
					{
						"id":    44,
						"label": "Full Stack",
						"value": "full stack",
					},
					{
						"id":    72,
						"label": "Full Stack (Java+React.js)",
						"value": "full stack (java+react.js)",
					},
					{
						"id":    79,
						"label": "Full Stack (PHP + JS)",
						"value": "full stack (php + js)",
					},
					{
						"id":    60,
						"label": "Full stack (React\\PHP)",
						"value": "full stack (react\\php)",
					},
					{
						"id":    75,
						"label": "Full-Stack (Angular\\Node)",
						"value": "full-stack (angular\\node)",
					},
					{
						"id":    59,
						"label": "Full-Stack (Node\\React)",
						"value": "full-stack (node\\react)",
					},
					{
						"id":    73,
						"label": "Fullstack (Java+Angular)",
						"value": "fullstack (java+angular)",
					},
					{
						"id":    49,
						"label": "Game Designer",
						"value": "game designer",
					},
					{
						"id":    56,
						"label": "Game Developer (JavaScript)",
						"value": "game developer (javascript)",
					},
					{
						"id":    65,
						"label": "HR",
						"value": "hr",
					},
					{
						"id":    32,
						"label": "Html5",
						"value": "html5",
					},
					{
						"id":    7,
						"label": "iOS",
						"value": "ios",
					},
					{
						"id":    70,
						"label": "IT Recruiter",
						"value": "it recruiter",
					},
					{
						"id":    1,
						"label": "Java",
						"value": "java",
					},
					{
						"id":    86,
						"label": "Java Script",
						"value": "java script",
					},
					{
						"id":    94,
						"label": "JS",
						"value": "js",
					},
					{
						"id":    11,
						"label": "Magento",
						"value": "magento",
					},
					{
						"id":    46,
						"label": "Marketing",
						"value": "marketing",
					},
					{
						"id":    31,
						"label": "Mobile",
						"value": "mobile",
					},
					{
						"id":    76,
						"label": "Node.js",
						"value": "node.js",
					},
					{
						"id":    84,
						"label": "Office Administrator",
						"value": "office administrator",
					},
					{
						"id":    62,
						"label": "Perl",
						"value": "perl",
					},
					{
						"id":    3,
						"label": "PHP",
						"value": "php",
					},
					{
						"id":    13,
						"label": "PM",
						"value": "pm",
					},
					{
						"id":    66,
						"label": "PostgreSQL",
						"value": "postgresql",
					},
					{
						"id":    41,
						"label": "Product Owner",
						"value": "product owner",
					},
					{
						"id":    85,
						"label": "Project Manager",
						"value": "project manager",
					},
					{
						"id":    9,
						"label": "Python",
						"value": "python",
					},
					{
						"id":    8,
						"label": "QA",
						"value": "qa",
					},
					{
						"id":    61,
						"label": "QA (Automation_Java)",
						"value": "qa (automation_java)",
					},
					{
						"id":    68,
						"label": "QA Automation",
						"value": "qa automation",
					},
					{
						"id":    92,
						"label": "QA Manual",
						"value": "qa manual",
					},
					{
						"id":    33,
						"label": "R&D",
						"value": "r&d",
					},
					{
						"id":    57,
						"label": "React Native",
						"value": "react native",
					},
					{
						"id":    95,
						"label": "Recommendation",
						"value": "recommendation",
					},
					{
						"id":    4,
						"label": "Ruby",
						"value": "ruby",
					},
					{
						"id":    67,
						"label": "Salesforce",
						"value": "salesforce",
					},
					{
						"id":    90,
						"label": "Scala",
						"value": "scala",
					},
					{
						"id":    74,
						"label": "SEO",
						"value": "seo",
					},
					{
						"id":    93,
						"label": "Solidity Developer",
						"value": "solidity developer",
					},
					{
						"id":    26,
						"label": "Solutions Architect",
						"value": "solutions architect",
					},
					{
						"id":    40,
						"label": "Support engineer",
						"value": "support engineer",
					},
					{
						"id":    82,
						"label": "Systems Engineer",
						"value": "systems engineer",
					},
					{
						"id":    51,
						"label": "Tech Artist",
						"value": "tech artist",
					},
					{
						"id":    6,
						"label": "UI/UX",
						"value": "ui/ux",
					},
					{
						"id":    10,
						"label": "Unity 3d",
						"value": "unity 3d",
					},
					{
						"id":    54,
						"label": "VR programmer",
						"value": "vr programmer",
					},
					{
						"id":    55,
						"label": "VR/AR Art director",
						"value": "vr/ar art director",
					},
					{
						"id":    43,
						"label": "Web UI",
						"value": "web ui",
					},
					{
						"id":    69,
						"label": "WMS Implementation & Support Consultants",
						"value": "wms implementation & support consultants",
					},
					{
						"id":    34,
						"label": "Xamarin",
						"value": "xamarin",
					},
				},
				"recruiters": []map[string]interface{}{
					{
						"id":    21,
						"label": "Dima",
						"value": "dima",
					},
					{
						"id":    29,
						"label": "Inna Parhomchuk",
						"value": "inna parhomchuk",
					},
					{
						"id":    30,
						"label": "Hanna Skorokhod",
						"value": "hanna skorokhod",
					},
					{
						"id":    31,
						"label": "Ludmila Samchenko",
						"value": "ludmila samchenko",
					},
					{
						"id":    32,
						"label": "Yana Miron",
						"value": "yana miron",
					},
					{
						"id":    33,
						"label": "Alona Rydun",
						"value": "alona rydun",
					},
					{
						"id":    35,
						"label": "Irina Krukovich",
						"value": "irina krukovich",
					},
					{
						"id":    38,
						"label": "Test",
						"value": "test",
					},
					{
						"id":    39,
						"label": "Kate",
						"value": "kate",
					},
					{
						"id":    40,
						"label": "Anastasiia",
						"value": "anastasiia",
					},
					{
						"id":    41,
						"label": "Kristina Miniailo",
						"value": "kristina miniailo",
					},
				},
				"reject_reasons": []map[string]interface{}{
					{
						"id":    5,
						"label": "closed to offers",
						"value": "closed to offers",
					},
					{
						"id":    11,
						"label": "does not fit",
						"value": "does not fit",
					},
					{
						"id":    8,
						"label": "does not like the project",
						"value": "does not like the project",
					},
					{
						"id":    6,
						"label": "low salary rate",
						"value": "low salary rate",
					},
					{
						"id":    10,
						"label": "remote only",
						"value": "remote only",
					},
					{
						"id":    9,
						"label": "terms don’t fit",
						"value": "terms don’t fit",
					},
					{
						"id":    7,
						"label": "was contacted earlier",
						"value": "was contacted earlier",
					},
				},
				"reject_tag": []map[string]interface{}{
					{
						"id":    3,
						"label": "reject",
						"value": "reject",
					},
				},
				"seniorities": []map[string]interface{}{
					{
						"id":    5,
						"label": "Architect",
						"value": "architect",
					},
					{
						"id":    1,
						"label": "Jun",
						"value": "jun",
					},
					{
						"id":    6,
						"label": "Jun-Mid",
						"value": "jun-mid",
					},
					{
						"id":    4,
						"label": "Lead",
						"value": "lead",
					},
					{
						"id":    2,
						"label": "Mid",
						"value": "mid",
					},
					{
						"id":    7,
						"label": "Mid-Sen",
						"value": "mid-sen",
					},
					{
						"id":    3,
						"label": "Sen",
						"value": "sen",
					},
					{
						"id":    8,
						"label": "Sen-Lead",
						"value": "sen-lead",
					},
				},
				"tags": []map[string]interface{}{
					{
						"id":    1,
						"label": "first contact",
						"value": "first contact",
					},
					{
						"id":    2,
						"label": "interested",
						"value": "interested",
					},
					{
						"id":    4,
						"label": "no answer\r\n",
						"value": "no answer\r\n",
					},
					{
						"id":    3,
						"label": "reject",
						"value": "reject",
					},
				},
				"vacancyStatus": []map[string]interface{}{
					{
						"id":    0,
						"label": "Hot",
						"value": "hot",
					},
					{
						"id":    1,
						"label": "Open",
						"value": "open",
					},
					{
						"id":    2,
						"label": "In Progress",
						"value": "in progress",
					},
					{
						"id":    3,
						"label": "Closed",
						"value": "closed",
					},
				},
			},
			"token_type": "Bearer",
			"user": map[string]interface{}{
				"active":           nil,
				"checklist":        nil,
				"company_id":       nil,
				"count_platform":   nil,
				"email":            "julia@uppeople.co",
				"id":               float64(27),
				"image":            nil,
				"name":             "Julia",
				"password":         "$2y$10$zyDhqar3N./1dBrLg3VE6elvuN8xK2C0YAOcJ5.aXA.EWuDSM9m36",
				"platform":         nil,
				"role_id":          1,
				"tel":              nil,
				"user_freelancers": nil,
			},
		}
		u := &auth.User{
			Companies:   make(map[int32]map[string]string),
			Host:        hosts[1],
			TokenOld:    v["access_token"].(string),
			UsersFields: &db.UsersFields{Id: int32(v["user"].(map[string]interface{})["id"].(float64))},
		}

		var err error
		u.Token, err = auth.Bearer.NewToken(u)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusNonAuthoritativeInfo)
			return nil, errors.Wrap(err, "Bearer.NewToken")
		}

		v["access_token"] = u.Token

		return v, nil
	}

	req := &fasthttp.Request{}
	ctx.Request.CopyTo(req)

	for _, host := range hosts {
		resp := &fasthttp.Response{}

		err := doRequest(req, resp, host)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode() == fasthttp.StatusOK {
			var v map[string]interface{}
			err := jsoniter.Unmarshal(resp.Body(), &v)
			u := &auth.User{
				Companies:   make(map[int32]map[string]string),
				Host:        host,
				TokenOld:    v["access_token"].(string),
				UsersFields: &db.UsersFields{Id: int32(v["user"].(map[string]interface{})["id"].(float64))},
			}

			u.Token, err = auth.Bearer.NewToken(u)
			if err != nil {
				ctx.SetStatusCode(fasthttp.StatusNonAuthoritativeInfo)
				return nil, errors.Wrap(err, "Bearer.NewToken")
			}

			v["access_token"] = u.Token

			return v, nil
		}
	}

	return nil, nil
}
