CREATE OR REPLACE FUNCTION get_log(_id integer, _isCand bool)
    RETURNS table(
                     logId integer,
                     text text,
                     date timestamp with time zone
                 )
    LANGUAGE plpgsql
AS
$$
BEGIN
    return query
    select logs.id as logId,
       CONCAT('Рекрутер ',
              users.name,
              log_actions.text_before_cand,
              CASE
                  WHEN logs.candidate_id > 0 THEN CONCAT(log_actions.for_candidate, can.name)
                  ELSE CONCAT(log_actions.for_company, companies.name) END,
              log_actions.text_after_cand,
              CASE
                  WHEN log_actions.name = 'CODE_SEND_CV' or log_actions.name = 'CODE_APPOINT_INTERVIEW'
                      THEN Format(' на вакансию %s, %s %s  в компанию %L', platforms.name, seniorities.name,
                                  CASE WHEN vacancies.name is not null
                                           THEN Format(' (%s)', vacancies.name)
                                       ELSE '' END,
                                  companies.name)

                  ELSE '' END,
              CASE WHEN log_actions.is_insert_text THEN
                       CASE WHEN (log_actions.name = 'CODE_LOG_UPDATE'
                           AND logs.changed is not null AND jsonb_typeof(logs.changed) = 'object')
                                THEN

                                (SELECT string_agg(Format(
                                    ' %s%s %s',
                                        CASE
                                            WHEN jst.key::text = 'platforms' THEN 'platforms'
                                            WHEN jst.key::text = 'platform_id' THEN 'platforms'
                                            WHEN jst.key::text = 'seniority_id' THEN 'seniority'
                                            WHEN jst.key::text = 'id_languages' THEN 'language_level'
                                            WHEN jst.key::text = 'tag_id' THEN 'tag/reject_reason'
                                            WHEN jst.key::text = 'vacancy_id' THEN 'по вакансии'
                                            WHEN jst.key::text = 'vacancies' THEN 'по вакансиям'
                                            WHEN jst.key::text = 'status_for_vac' THEN 'cтатус по вакансии'
                                            WHEN jst.key::text = 'contact_id' THEN 'контакт'
                                            WHEN jst.key::text = 'text' THEN ''
                                            ELSE jst.key::text END,

                                        CASE
                                            WHEN jst.key::text = 'text' THEN ''
                                            ELSE ': ' END,
                                        CASE
                                            WHEN jst.key::text = 'platforms'
                                                THEN (select string_agg(ps.name, ', ') from platforms ps where ps.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'platform_id'
                                                THEN (select string_agg(ps.name, ', ') from platforms ps where ps.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'seniority_id'
                                                THEN (select string_agg(ss.name, ', ') from seniorities ss where ss.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'id_languages'
                                                THEN (select string_agg(ls.name, ', ') from languages ls where ls.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'tag_id'
                                                THEN (select string_agg(ts.name, ', ') from tags ts where ts.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'status_for_vac'
                                                THEN (select string_agg(sv.status, ', ') from status_for_vacs sv where sv.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'contact_id'
                                                THEN (select string_agg(cs.name, ', ') from contacts cs where cs.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN (jst.key::text = 'vacancy_id' OR jst.key::text = 'vacancies')
                                                THEN (select string_agg(
                                                        Format('%s, %s %s в компании %L',
                                                        (select pss.name from platforms pss WHERE vs.platform_id = pss.id),
                                                        (select sss.name FROM seniorities sss WHERE vs.seniority_id = sss.id),
                                                        CASE WHEN vs.name is not null THEN Format(' (%s)', vs.name) ELSE '' END,
                                                        (select css.name FROM companies css WHERE vs.company_id = css.id)
                                                        ), ', ')
                                                      FROM vacancies vs
                                                      where vs.id = ANY(jsonb_array_castint(jst.value)))


                                            ELSE jst.value::text END
                                    ), ', ')

                                 FROM jsonb_each(logs.changed::jsonb) jst)

                            ELSE CONCAT(' ', logs.text) END
                   ELSE '' END
           ) as text,
       logs.create_at as date

from logs left Join companies on (logs.company_id = companies.id)
    left join vacancies ON (logs.vacancy_id = vacancies.id)
    join users ON (logs.user_id = users.id)
    join candidates can ON (logs.candidate_id = can.id)
    left Join platforms ON (vacancies.platform_id = platforms.id)
    left Join seniorities ON (vacancies.seniority_id = seniorities.id)
    left Join log_actions ON (logs.action_code = log_actions.id)
where (_isCand AND logs.candidate_id = _id) or (NOT _isCand AND logs.company_id = _id)
order by logs.create_at DESC
;
END;
$$;