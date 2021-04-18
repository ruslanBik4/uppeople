create table statuses
(
    id serial not null,
    status character varying not null,
        PRIMARY KEY (id)
);
COMMENT ON TABLE statuses IS 'statuses of vacancies';
-- example comment with dataJSON
COMMENT ON COLUMN statuses.status IS 'full name';
-- examply index
create unique index statuses_name_idx
    ON statuses (status);
