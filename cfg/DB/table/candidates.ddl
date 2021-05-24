create table candidates
(
    id serial not null,
    platform_id integer not null default 0,
    name character varying not null,
    salary integer not null default 0,
    email character varying not null default '',
    phone character varying not null default '',
    skype character varying not null default '',
    link character varying not null default '',
    linkedin character varying default '',
    str_companies character varying default '',
    status character varying not null default '',
    tag_id integer not null default 1,
    comments text not null default '',
    date timestamp with time zone not null default CURRENT_TIMESTAMP,
    recruter_id integer not null default 1,
    text_rezume text not null default '',
    sfera character varying not null default '',
    experience character varying not null default '',
    education character varying not null default '',
    language character varying not null default '',
    zapoln_profile integer,
    file character varying not null default '',
    avatar character varying not null default '',
    seniority_id integer not null default 1,
    date_follow_up date,
    vacancies integer[],
        PRIMARY KEY (id)
);
COMMENT ON TABLE candidates IS 'list of candidates';

create unique index candidates_name_uindex
    on candidates (name);

create unique index candidates_email_uindex
    on candidates (email)
    where (((email)::text > ''::text) AND (email IS NOT NULL));

create unique index candidates_mobile_uindex
    on candidates (phone)
    where (((phone)::text > ''::text) AND (phone IS NOT NULL));

create unique index candidates_linkedin_uindex
    on candidates (linkedin)
    where (((linkedin)::text > ''::text) AND (linkedin IS NOT NULL));

alter table candidates
    add constraint candidates_seniorities_id_fk
        foreign key (seniority_id) references seniorities
            on update cascade on delete set default;

alter table candidates
    add constraint candidates_tags_id_fk
        foreign key (tag_id) references tags
            on update cascade on delete set default
