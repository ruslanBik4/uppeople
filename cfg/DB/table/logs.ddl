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

update logs
set changed=(
    CASE WHEN text LIKE '%=https:%' OR text LIKE '%comments=%' OR text LIKE '%details%' THEN null
         WHEN text LIKE '%=%' THEN
             (SELECT CONCAT('{',
                            array_to_string(
                                    array_agg(
                                            CONCAT('"',
                                                   REPLACE(
                                                           CASE
                                                               WHEN val LIKE '%[%]%'
                                                                   THEN REPLACE(val, ' ', ',')
                                                               WHEN val LIKE '%<p>%</p>%' THEN CONCAT(regexp_replace(
                                                                                                              REPLACE(REPLACE(REPLACE(val, '<p>', ''), '</p>', ''), '=', '="'),
                                                                                                              E'[\\n\\r\\f\\u000B\\u0085\\u2028\\u2029]+', '', 'g'), '"')
                                                               WHEN substring(val from '(?<==)[^\]]+') ~ '^0+\d+$' THEN CONCAT(REPLACE(val, '=', '="'), '"')
                                                               WHEN substring(val from '(?<==)[^\]]+') ~ '^\d+(\.\d)?$' THEN val
                                                               ELSE CONCAT(REPLACE(val, '=', '="'), '"') END
                                                       , '=', '":')
                                                )
                                        ), ','), '}')::jsonb
              FROM
                  regexp_split_to_table(text, E', ') as x(val))
         WHEN text LIKE '{%}' THEN text::jsonb
         ELSE to_jsonb(text) END)
where logs.text IS NOT NULL AND logs.changed IS NULL;

