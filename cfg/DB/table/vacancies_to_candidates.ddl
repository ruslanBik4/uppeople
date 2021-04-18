create table vacancies_to_candidates
(
    candidate_id integer not null,
    company_id integer,
--         constraint vacancies_to_candidates_companies_id_fk
--             references companies
--             on update cascade on delete cascade,
    vacancy_id integer not null,
--         constraint vacancies_to_candidates_vacancies_id_fk
--             references vacancies
--             on update cascade on delete cascade,
    status integer,
    user_id integer,
    date_create timestamp default CURRENT_TIMESTAMP,
    date_last_change date,
    rej_text text not null default ''::text,
    rating integer,
    notice text,
    constraint vacancies_to_candidates_pk
        primary key (candidate_id, vacancy_id)
);

COMMENT ON TABLE vacancies_to_candidates IS 'statuses of vacancies candidate workflow';

