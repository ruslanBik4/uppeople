CREATE OR REPLACE FUNCTION get_log_test(Id integer, isCand bool)
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
              CASE WHEN log_actions.is_insert_text = true THEN CONCAT(' ', logs.text) ELSE '' END
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