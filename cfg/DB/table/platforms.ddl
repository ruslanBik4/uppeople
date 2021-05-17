create table platforms
(
    id   serial not null,
    name varchar(255) not null,
    primary key (id)
);

COMMENT ON TABLE platforms IS 'platforms of candidates';
-- example comment with dataJSON
COMMENT ON COLUMN platforms.name IS 'full name';
-- examply index
create unique index platforms_name_idx
    ON platforms (name);

