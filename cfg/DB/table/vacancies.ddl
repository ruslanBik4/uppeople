create table vacancies
(
    id serial not null,
    company_id bigint,
    platform_id bigint,
    user_ids integer[] not null,
    name varchar(255) default NULL::character varying,
    description text,
    details text,
    link text default ''::text not null,
    file varchar(255) default NULL::character varying,
    date_create timestamp default CURRENT_TIMESTAMP not null,
    ord bigint,
    status bigint,
    seniority_id bigint not null,
    salary bigint not null,
    location_id bigint,
    primary key (id)
);

-- todo: https://commitfest.postgresql.org/17/1252/


create index idx_17194_vacancies_seniority_id_foreign
    on vacancies (seniority_id);

create index idx_17194_vacancies_status_foreign
    on vacancies (status);

create index idx_17194_vacancies_location_id_foreign
    on vacancies (location_id);

create index idx_17194_vacancies_platform_id_foreign
    on vacancies (platform_id);

create index idx_17194_vacancies_company_id_foreign
    on vacancies (company_id);

