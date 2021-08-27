CREATE TABLE companies
(
    id   serial,
        name character varying not null,
        send_details text,
        interview_detail text,
        cooperation text,
        contact text,
        about text,
        map text,
        phone varchar(255) default NULL::character varying,
        email varchar(255) default NULL::character varying,
        skype varchar(255) default NULL::character varying,
        logo varchar(255) default NULL::character varying,
        address text,
        email_template text,
        manager_id bigint
        );


create unique index companies_name_uindex
    on companies (name);

