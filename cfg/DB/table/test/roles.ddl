CREATE TABLE roles
(
    id     SERIAL            NOT NULL
        CONSTRAINT roles_pkey PRIMARY KEY,
    name   character varying NOT NULL,
    grants JSON,
    views  JSON
);
COMMENT ON TABLE roles IS 'List of privilege of users';

create unique index if not exists roles_name_idx
    ON roles (name);

INSERT INTO roles (name, grants, views)
VALUES ('admin', json('{"all":true}'), json('{"all":true}')),
       ('manager', json('{"grant":"manager"}'), json('{}')),
       ('user', json('{}'), json('{}'))
