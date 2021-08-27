CREATE TABLE contacts
(
    id              serial,
    company_id      integer,
    name            character varying not null,
    email           character varying not null,
    phone           character varying,
    skype           character varying,
    default_contact integer,
    all_platforms   integer,
    not_visible     integer,
    PRIMARY KEY (id)
);

COMMENT ON TABLE contacts IS '';
-- example comment with dataJSON
COMMENT ON COLUMN contacts.name IS 'full name {"pattern": "name"}';
-- examply index
create unique index contacts_email_idx
    ON contacts (email, company_id);


alter table contacts
    add constraint contacts_company_id_fk
        foreign key (company_id) references companies
            on update cascade on delete cascade;
