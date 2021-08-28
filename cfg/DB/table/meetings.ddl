CREATE TABLE meetings
(
    id           serial,
    candidate_id integer,
    company_id   integer,
    vacancy_id   integer,
    title        character varying,
    date         timestamp with time zone,
    user_id      integer,
    color        character varying,
    type         character varying
);

COMMENT ON TABLE meetings IS 'Calendar of meetings of group';

alter table meetings
    add constraint meetings_company_id_fk
        foreign key (company_id) references companies
            on update cascade on delete set default;

alter table meetings
    add constraint meetings_vacancy_id_fk
        foreign key (vacancy_id) references vacancies
            on update cascade on delete set default;

alter table meetings
    add constraint meetings_user_id_fk
        foreign key (user_id) references public.users
            on update cascade on delete set default;

alter table meetings
    add constraint meetings_candidate_id_fk
        foreign key (candidate_id) references candidates
            on update cascade on delete set default;