// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/valyala/fasthttp"
)

func HandleDashBoard(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	return DB.Conn.SelectToMap(ctx,
		`with vac as (
    select *
    from vacancies
    WHERE status =ANY(ARRAY[1,0])
),
      intRows as (
          select *
          from int_rev_candidates
            where status =ANY(ARRAY [2,9])
      )
select json_build_object('countVac', count(id), 'countCom', count(distinct company_id)) "countVacanciesOpenAndHot",
       (select json_object_agg(name, cId)
           from (select name, count(v.id) cId
                 from platforms p join vac v on p.id = v.platform_id
           group by name) j
           ) "countVacanciesOpenAndHotForPlatform",
       (select json_object_agg(name, j.obj)
           from (select c.name, json_build_object('countAll', count(i.id), 'Review',
               count(i.id) FILTER ( WHERE s.id = 1 )) obj
                 from intRows i join vacancies v on (i.vacancy_id = v.id)
                     join platforms p on p.id = v.platform_id
                     join companies c on i.company_id = c.id
                     join status_for_vacs s on i.status = s.id
           group by c.name) j
           ) "countRevInterComp",
       (select json_object_agg(vacancy, j.obj)
        from (select vacancy, json_build_object('countAll', count(i.id),
                                        'Review',
                                       count(i.id) FILTER ( WHERE s.id = 9 )) obj
              from intRows i join vacancies v on (i.vacancy_id = v.id)
                             join platforms p on p.id = v.platform_id
                             join companies c on i.company_id = c.id
                            join seniorities on v.seniority_id = seniorities.id
                             join status_for_vacs s on i.status = s.id
                             JOIN lateral CONCAT_WS(' - ',c.name, p.name, seniorities.name) vacancy on true
              group by vacancy) j
       ) "countRevInterVac",
       (select json_agg(j.obj)
        from (select json_build_object('id', v.id,
                                        'vacancy', vacancy) obj
              from vac v join platforms p on p.id = v.platform_id
                             join companies c on v.company_id = c.id
                            join seniorities on v.seniority_id = seniorities.id
                             JOIN lateral CONCAT_WS(' - ',c.name, p.name, seniorities.name) vacancy on true
              where age(v.date_create) < interval '10 day') j
       ) "listNotSendVacancies",
      (select json_object_agg(role, uId)
      from (select roles.nazva_en as role, count(users.id) uId
            from users join roles  on users.role_id = roles.id
          group by 1) j) "countUsersAndFreelancer",
    (select json_object_agg(status, j.obj)
        from (select s.status, count(i.id)  obj
              from int_rev_candidates i
                join status_for_vacs s on i.status = s.id
              where age(i.date) < interval '1 month'
              group by 1) j
       ) "countReviewAndInterviewCandidatesOnVacancies",
    (select json_object_agg(status, j.obj)
        from (select s.status, json_agg( json_build_object( 'id', i.id, 'name', users.name, 'role', roles.nazva_en) ) obj
              from int_rev_candidates i
                join status_for_vacs s on i.status = s.id
                join users on i.user_id=users.id
                join roles on users.role_id = roles.id
              where age(i.date) < interval '1 month'
              group by 1) j
       ) "topUsers",
       (select json_build_object( 'allCount', count(users.id),
            'users', json_agg( json_build_object('name', users.name, 'id', users.id) ))
              from vacancies_to_candidates v
                       join users on v.user_id=users.id
                       join roles on users.role_id = roles.id
              where v.status = 5 and age(coalesce(v.date_last_change, v.date_create)) < interval '1 month'

       ) "offersCandidates"
from vac`)
}
