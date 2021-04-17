create table statuses
(
    id serial not null,
    status varchar(255) default NULL::character varying,
        PRIMARY KEY (id)
);
COMMENT ON TABLE statuses IS 'statuses of vacancies';
-- example comment with dataJSON
COMMENT ON COLUMN statuses.status IS 'full name {"pattern": "name"}';
-- examply index
create unique index statuses_name_idx
    ON statuses (status);
