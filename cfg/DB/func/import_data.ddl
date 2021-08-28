CREATE OR REPLACE procedure import_data()
LANGUAGE plpgsql
AS
$$
BEGIN

    truncate candidates cascade ;

--     EXPLAIN
    INSERT INTO candidates
        (select id,
                array[platform_id],
                name,
                (CASE WHEN salary ~ '\d+-\d+'
                    THEN (regexp_match(salary, '(\d+)-(\d+)'))[2]::integer
                    WHEN salary ~ '^\d{1,7}$' THEN salary::integer
                    WHEN salary ~ '^\d{8}$' THEN "left"( salary, 7)::integer
                    ELSE 0 END),
                coalesce(email, ''),
                coalesce(phone,''),
                coalesce(skype,''),
                coalesce(link,''),
                linkedin,
                coalesce(status,''),
                tag_id,
                coalesce(comments,''),
                date,
                recruter_id,
                coalesce(text_rezume, '') as cv,
                coalesce(experience,''),
                coalesce(education,''),
                (select l.id from public.languages l where l.name = coalesce(language,'Unknown')),
                coalesce(file,''),
                avatar::bytea,
                coalesce(seniority_id,1),
                date_follow_up,
                array(select v.id from vacancies_to_candidates_tmp v where candidate_id = id AND status = 1)
         from candidates_tmp c
         where c.name >'' AND not exists(select null from candidates_tmp t
                                         where c.id < t.id AND
                                             ((c.name = t.name)
                                                 or (c.linkedin > '' and c.linkedin = t.linkedin)
                                                 or (c.email > '' AND c.email = t.email)
                                                 or (c.phone > '' AND c.phone = t.phone))
             ))
    on conflict do nothing ;
    PERFORM setval('candidates_id_seq'::regclass, (select max(id) from candidates));

    truncate table vacancies cascade;

--     EXPLAIN
    INSERT INTO vacancies
        (select id,
                company_id,
                platform_id,
                format('{%s}', user_ids)::integer[],
                name,
                description,
                details,
                coalesce(link,''),
                file,
                date_create,
                ord,
                status,
                seniority_id,
                coalesce(salary,0),
                coalesce(location_id,0)
         from vacancies_tmp)
    on conflict do nothing ;

    PERFORM setval('vacancies_id_seq'::regclass, (select max(id) from vacancies));


    insert into vacancies_to_candidates
        (select

                candidate_id,
                company_id,
                vacancy_id,
                status,
                user_id,
                date_create,
                date_last_change,
                coalesce(rej_text, ''),
                rating,
                notice
         from vacancies_to_candidates_tmp v
        where v.candidate_id in (select id from candidates)
              AND not exists(select null from vacancies_to_candidates_tmp t
                         where v.id < t.id AND
                             (v.vacancy_id = t.vacancy_id)
                                 AND (v.candidate_id = t.candidate_id)))
    on conflict do nothing ;

    PERFORM setval('companies_id_seq'::regclass, (select max(id) from companies));
    PERFORM setval('comments_for_companies_id_seq'::regclass, (select max(id) from comments_for_companies));
    PERFORM setval('comments_for_candidates_id_seq'::regclass, (select max(id) from comments_for_candidates));
--     PERFORM setval('candidates_to_companies_id_seq'::regclass, (select max(id) from candidates_to_companies));
    PERFORM setval('contacts_id_seq'::regclass, (select max(id) from contacts));

END;

$$;
