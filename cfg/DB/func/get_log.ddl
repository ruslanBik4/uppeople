CREATE OR REPLACE FUNCTION get_log(Id integer, isCand bool)
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
              CASE WHEN log_actions.is_insert_text = true THEN
                       CASE WHEN (log_actions.name = 'CODE_LOG_UPDATE'
                           AND logs.text LIKE '{%}'
                           AND logs.text::json is not null)
                                THEN

                                (SELECT DISTINCT array_to_string(array_agg(CONCAT(
                                        CASE
                                            WHEN jst.key::text = 'platform_id' THEN 'platforms'
                                            WHEN jst.key::text = 'seniority_id' THEN 'seniority'
                                            WHEN jst.key::text = 'id_languages' THEN 'language_level'
                                            WHEN jst.key::text = 'tag_id' THEN 'tag/reject_reason'
                                            WHEN jst.key::text = 'vacancy_id' THEN 'по вакансии'
                                            ELSE jst.key::text END,
                                        '=',
                                        CASE
                                            WHEN jst.key = 'platform_id' THEN (select array_to_string(array_agg(name), ', ') from platforms where id = ANY(json_array_castint(jst.value)))
                                            WHEN jst.key = 'seniority_id' THEN (select array_to_string(array_agg(name), ', ') from seniorities where id = ANY(json_array_castint(jst.value)))
                                            WHEN jst.key = 'id_languages' THEN (select array_to_string(array_agg(name), ', ') from languages where id = ANY(json_array_castint(jst.value)))
                                            WHEN jst.key = 'tag_id' THEN (select array_to_string(array_agg(name), ', ') from tags where id = ANY(json_array_castint(jst.value)))
                                            WHEN jst.key = 'vacancy_id'
                                                THEN (select CONCAT(
                                                                     platforms.name,
                                                                     ', ',
                                                                     seniorities.name,
                                                                     CASE WHEN vacancies.name is not null THEN CONCAT(' (', vacancies.name, ')') ELSE '' END,
                                                                     ' в компании ',
                                                                     companies.name)
                                                      FROM vacancies
                                                               left join companies on (vacancies.company_id = companies.id)
                                                               left Join platforms ON (vacancies.platform_id = platforms.id)
                                                               left Join seniorities ON (vacancies.seniority_id = seniorities.id)
                                                      where vacancies.id = jst.value::text::integer)
                                            ELSE jst.value::text END
                                            )), ', ')

                                 FROM json_each(logs.text::json) jst)

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
where (logs.candidate_id = $1 AND $2 = true) or (logs.company_id = $1 AND $2 = false)
order by logs.create_at DESC
;
END;
$$;