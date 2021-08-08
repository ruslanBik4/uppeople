create table logs
(
    id serial not null,
    user_id integer not null default 1,
    candidate_id integer,
    company_id integer,
    vacancy_id integer,
    text text,
    changed jsonb,
    action_code integer not null,
    date_create date not null,
    create_at timestamp with time zone not null default CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

alter table logs
    add constraint logs_candidates_id_fk
        foreign key (candidate_id) references candidates
            on update cascade on delete cascade;
alter table logs
    add constraint logs_users_id_fk
        foreign key (user_id) references users
            on update cascade on delete set default;
-- alter table logs
--     add constraint logs_company_id_fk
--         foreign key (company_id) references companies
--             on update cascade on delete set default;
-- alter table logs
--     add constraint logs_vacancy_id_fk
--         foreign key (vacancy_id) references vacancies
--             on update cascade on delete cascade;

update logs
set action_code=(select id from log_actions where name = 'CODE_APPOINT_INTERVIEW')
where logs.text LIKE '% назначил %';

update logs
set action_code=(select id from log_actions where name = 'CODE_SEND_CV')
where logs.text LIKE '% CV %';

CALL update_log();

