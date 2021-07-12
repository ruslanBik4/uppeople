// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

// names of system environment variables
const (
	CFG_PATH    = "configPath"
	WEB_PATH    = "webPath"
	SYSTEM_PATH = "systemPath"
)

const LOG_VIEW = "select * from get_log($1, $2)"

const (
	EMAIL_TEXT = `<p><span style="font-size: 14px;">Colleagues,
please review the candidacy of %s for the position of %s </span></p>
<p>CV:%s</p>
<p>Experience:%s</p>
<p>English level:%s</p>
<p>Salary expectations:$%v</p>
<p><br>Will be appreciate for quick feedback.</p>
<p><br><br></p>
<p>Best regards,
UPPeople team.</p>
<p><span style="font-size: 14px;">Добрый день,
рассмотрите, пожалуйста, кандидата %[1]s на позицию  %s </span></p>
<p>CV:%s</p>
<p>Опыт:%s</p>
<p>Уровень английского:%s</p>
<p>Ожидания по заработной плате:$%v</p>
<p><br>Будем благодарны за фидбек.</p>
<p><br><br></p>
<p>С наилучшим пожеланиями,
команда UPPeople.</p>
<p>&nbsp;
<a href="http://my.uppeople.co/" target="_self"><span style="color: blue;font-size: 16px;font-family: Journal, serif;"@"UPpeople" Recruiting agency</span></a><span style="font-size: 16px;"> </span></p>
`
	EMAIL_INVITE_TEXT = `<p><span style="font-size: 14px;">Interview with %s %s  scheduled on %s %s</span></p>
<p>%s</p>
<p><br>Add to google calendar: {link}</p>
<p><br></p>
<p><span style="color: rgb(35,40,44);background-color: rgb(255,255,255);font-size: 14px;">
<a href="http://my.uppeople.co/" target="_self"><span style="color: blue;font-size: 16px;font-family: Journal, serif;"@"UPpeople" Recruiting agency</span></a><span style="font-size: 16px;"> </span></p>
 <br><br>`
)

const SQL_VIEW_CANDIDATE_VACANCIES = `select v.id,
		j.name, 
		j.name as label, 
		LOWER(j.name) as value, 
		user_ids, 
		platform_id,
        coalesce( (select u.name from users u where u.id = vc.user_id), '') as recruiter,
		CONCAT(platforms.name, ' ("', 
			(select name from seniorities where id=seniority_id), '")') as platform,
		companies, sv.id as status_id, v.company_id, sv.status, salary, 
		coalesce(vc.date_last_change, vc.date_create) as date_last_change, vc.rej_text, sv.color
FROM vacancies v JOIN companies on (v.company_id=companies.id)
	JOIN vacancies_to_candidates vc on (v.id = vc.vacancy_id )
	JOIN platforms ON (v.platform_id = platforms.id)
	JOIN status_for_vacs sv on (vc.status = sv.id)
    JOIN LATERAL (select concat(companies.name, ' ("', platforms.name, '")') as name) j on true
	WHERE vc.candidate_id=$1 AND vc.status!=1
    order by date_last_change desc
`

const SEND_CV_COMPANIES_SQL = `SELECT id as comp_id, c.name, send_details,
  (select json_agg(json_build_object('email', t.email, 'id',t.id, 'all_platforms', t.all_platforms,
           'platform_id', cp.platform_id, 'name', t.name))
             from contacts t left join contacts_to_platforms cp on t.id=cp.contact_id
           WHERE t.company_id = c.id AND (all_platforms=1 OR platform_id=ANY($1))) as contacts,
  (select json_agg(json_build_object('id', v.id,
		   'platform', (select p.name  from platforms p where p.id = v.platform_id),
		   'location', (select l.name   from location_for_vacancies l where v.location_id = l.id),
           'seniority', (select s.name from seniorities s where s.id = v.seniority_id),
           'salary', v.salary, 
			'name', v.name, 
			'user_ids', v.user_ids)) as vacancy
	from vacancies v
	where c.id = v.company_id and status <= 1 and platform_id=ANY($1))
FROM companies c
WHERE c.id in (select v.company_id from vacancies v
    where status <= 1 and platform_id=ANY($1) and $2 = ANY(user_ids))`
