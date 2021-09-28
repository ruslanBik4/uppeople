create table platforms
(
    id   serial not null,
    name character varying not null,
    primary key (id)
);

COMMENT ON TABLE platforms IS 'platforms of candidates';
COMMENT ON COLUMN platforms.name IS 'Unique name';
create unique index platforms_name_idx
    ON platforms (name);