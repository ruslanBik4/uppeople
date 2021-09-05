CREATE TABLE int_rev_candidates
(
    candidate_id integer not null default 1,
    company_id   integer not null,
    vacancy_id   integer not null,
    status       integer not null,
    user_id      integer not null default 1,
    date         date not null,
    PRIMARY KEY (candidate_id, company_id, vacancy_id)
);
alter table int_rev_candidates
    add constraint int_rev_candidates_candidates_id_fk
        foreign key (candidate_id) references candidates
        on update cascade on delete cascade;
alter table int_rev_candidates
    add constraint int_rev_candidates_users_id_fk
        foreign key (user_id) references public.users
            on update cascade on delete set default;
alter table int_rev_candidates
    add constraint int_rev_candidates_status_fk
        foreign key (status) references public.status_for_vacs
            on update cascade on delete set default;
alter table int_rev_candidates
    add constraint int_rev_candidates_vacancies_id_fk
        foreign key (vacancy_id) references vacancies
            on update cascade on delete cascade;


COMMENT ON TABLE int_rev_candidates IS 'List candidates for campany on interview';

