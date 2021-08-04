CREATE OR REPLACE FUNCTION get_log(id integer, isCand bool)
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
                      THEN CONCAT(' на вакансию ', platforms.name, ', ', seniorities.name,
                                  CASE WHEN vacancies.name is not null
                                           THEN CONCAT(' (', vacancies.name, ')')
                                       ELSE '' END,
                                  ' в компанию ', companies.name)

                  ELSE '' END,
              CASE WHEN log_actions.is_insert_text THEN
                       CASE WHEN (log_actions.name = 'CODE_LOG_UPDATE'
                           AND logs.changed is not null AND jsonb_typeof(logs.changed) = 'object')
                                THEN

                                (SELECT DISTINCT string_agg(CONCAT(
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
                                            ELSE jst.key::text END,
                                        '=',
                                        CASE
                                            WHEN jst.key::text = 'platforms' THEN (select array_to_string(array_agg(name), ', ') from platforms ps where ps.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'platform_id' THEN (select array_to_string(array_agg(name), ', ') from platforms ps where ps.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'seniority_id' THEN (select array_to_string(array_agg(name), ', ') from seniorities ss where ss.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'id_languages' THEN (select array_to_string(array_agg(name), ', ') from languages ls where ls.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'tag_id' THEN (select array_to_string(array_agg(name), ', ') from tags ts where ts.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'status_for_vac' THEN (select array_to_string(array_agg(status), ', ') from status_for_vacs sv where sv.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN jst.key::text = 'contact_id' THEN (select array_to_string(array_agg(name), ', ') from contacts cs where cs.id = ANY(jsonb_array_castint(jst.value)))
                                            WHEN (jst.key::text = 'vacancy_id' OR jst.key::text = 'vacancies')
                                                THEN (select array_to_string(array_agg(CONCAT(
                                                    pss.name,
                                                    ', ',
                                                    sss.name,
                                                    CASE WHEN vs.name is not null THEN CONCAT(' (', vs.name, ')') ELSE '' END,
                                                    ' в компании ',
                                                    css.name)), ', ')
                                                      FROM vacancies vs
                                                               left join companies css on (vs.company_id = css.id)
                                                               left Join platforms pss ON (vs.platform_id = pss.id)
                                                               left Join seniorities sss ON (vs.seniority_id = sss.id)
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
where ($2 AND logs.candidate_id = $1) or (NOT $2 AND logs.company_id = $1)
order by logs.create_at DESC
;
END;
$$;