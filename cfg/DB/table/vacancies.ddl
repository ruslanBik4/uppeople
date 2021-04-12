create table vacancies
(
    id serial not null,
    company_id integer,
    platform_id integer,
    user_ids integer[] not null,
    name varchar(255) default NULL::character varying,
    description text,
    details text,
    link text not null default ''::text,
    file varchar(255) default NULL::character varying,
    date_create timestamp not null default CURRENT_TIMESTAMP,
    ord bigint,
    status bigint,
    seniority_id integer not null,
    salary bigint not null,
    location_id bigint,
    primary key (id)
);

-- todo: https://commitfest.postgresql.org/17/1252/
-- todo: add foreign keys

CREATE EXTENSION IF NOT EXISTS citext;

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

