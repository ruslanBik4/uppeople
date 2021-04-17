create table logs
(
    id serial not null,
    user_id integer,
    candidate_id integer,
    company_id integer,
    vacancy_id integer,
    text text not null,
    kod_deystviya integer not null,
    date_create date not null,
    create_at timestamp with time zone default CURRENT_TIMESTAMP not null,
    PRIMARY KEY (id)
);

