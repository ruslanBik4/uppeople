CREATE OR REPLACE procedure import_data()
LANGUAGE plpgsql
AS
$$
BEGIN

--     truncate candidates cascade ;

--     EXPLAIN
    INSERT INTO candidates
        (select id,
                platform_id,
                name,
                coalesce(salary,'0'),
                coalesce(email, ''),
                coalesce(phone,''),
                coalesce(skype,''),
                coalesce(link,''),
                linkedin,
                str_companies,
                coalesce(status,''),
                tag_id,
                coalesce(comments,''),
                date,
                recruter_id,
                coalesce(text_rezume, ''),
                coalesce(sfera,''),
                coalesce(experience,''),
                coalesce(education,''),
                coalesce(language,''),
                zapoln_profile,
                coalesce(file,''),
                coalesce(avatar,''),
                coalesce(seniority_id,0),
                date_follow_up
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
--     truncate table vacancies cascade;

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
                salary,
                location_id
         from vacancies_tmp )
    on conflict do nothing ;
    PERFORM setval('vacancies_id_seq'::regclass, (select max(id) from vacancies));


    insert into vacancies_to_candidates
        (select *
        from vacancies_to_candidates_tmp v
        where v.candidate_id in (select id from candidates)
              AND not exists(select null from vacancies_to_candidates_tmp t
                         where v.id < t.id AND
                             (v.vacancy_id = t.vacancy_id)
                                 AND (v.candidate_id = T.candidate_id)))
    on conflict do nothing ;
    PERFORM setval('vacancies_to_candidates_id_seq'::regclass, (select max(id) from vacancies_to_candidates));

    PERFORM setval('companies_id_seq'::regclass, (select max(id) from companies));
    PERFORM setval('comments_for_companies_id_seq'::regclass, (select max(id) from comments_for_companies));
    PERFORM setval('comments_for_candidates_id_seq'::regclass, (select max(id) from comments_for_candidates));
    PERFORM setval('candidates_to_companies_id_seq'::regclass, (select max(id) from candidates_to_companies));
    PERFORM setval('contacts_id_seq'::regclass, (select max(id) from contacts));

END;

$$;
