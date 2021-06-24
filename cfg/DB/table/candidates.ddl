create table candidates
(
    id serial not null,
    platforms integer[] not null,
    name character varying not null,
    salary integer not null default 0,
    email character varying not null default '',
    phone character varying not null default '',
    skype character varying not null default '',
    link character varying not null default '',
    linkedin character varying default '',
    status character varying not null default '',
    tag_id integer not null default 1,
    comments text not null default '',
    date timestamp with time zone not null default CURRENT_TIMESTAMP,
    recruter_id integer not null default 1,
    cv text not null default '',
    experience character varying not null default '',
    education character varying not null default '',
    language character varying not null default '',
    file character varying not null default '',
    avatar bytea not null default ''::bytea,
    seniority_id integer not null default 1,
    date_follow_up date,
    vacancies integer[],
    PRIMARY KEY (id)
);

COMMENT ON TABLE candidates IS 'list of candidates';
COMMENT ON COLUMN candidates.platforms IS 'Platform {"suggestions":"/api/main/returnOptionsForSelects", "multiple":true,"suggestions_params":{"name":"platforms"}}';
COMMENT ON COLUMN candidates.recruter_id IS 'Recruiter name {"suggestions":"/api/main/returnOptionsForSelects","suggestions_params":{"name":"recruiters"}}';
COMMENT ON COLUMN candidates.seniority_id IS 'Seniority {"suggestions":"/api/main/returnOptionsForSelects","suggestions_params":{"name":"seniorities"}}';
COMMENT ON COLUMN candidates.tag_id IS 'Tag {"suggestions":"/api/main/returnOptionsForSelects","suggestions_params":{"name":"tags"}}';
COMMENT ON COLUMN candidates.vacancies IS 'Vacancies {"suggestions":"/api/get_recruiter_vacancies","multiple":true}';

create unique index candidates_name_uindex
    on candidates (name);

create unique index candidates_email_uindex
    on candidates (email)
    where ((email)::text > ''::text);

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
