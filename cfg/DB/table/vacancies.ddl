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
    ord integer,
    status integer,
    seniority_id integer not null,
    salary integer not null,
    location_id integer,
    primary key (id)
);

alter table vacancies
    add constraint vacancies_seniorities_id_fk
        foreign key (seniority_id) references public.seniorities
            on update cascade on delete set default;

alter table vacancies
    add constraint vacancies_company_id_fk
        foreign key (company_id) references companies
            on update cascade on delete set default;

alter table vacancies
    add constraint vacancies_platform_id_fk
        foreign key (platform_id) references public.platforms
            on update cascade on delete set default;

alter table vacancies
    add constraint vacancies_location_id_fk
        foreign key (location_id) references public.location_for_vacancies
            on update cascade on delete set default;

alter table vacancies
    add constraint vacancies_status_fk
        foreign key (status) references public.status_for_vacs
            on update cascade on delete set default;


