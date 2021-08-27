CREATE TABLE comments_for_companies
(
    id          serial,
    user_id     integer,
    company_id  integer,
    comments    text,
    time_create character varying,
    PRIMARY KEY (id)
);

COMMENT ON TABLE comments_for_companies IS '';

alter table comments_for_companies
    add constraint comments_for_companies_company_id_fk
        foreign key (company_id) references companies
            on update cascade on delete cascade;

-- alter table comments_for_companies
--     add constraint comments_for_companies_user_id_fk
--         foreign key (user_id) references public.users
--             on update cascade on delete set default;