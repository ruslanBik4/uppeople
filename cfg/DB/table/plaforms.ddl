create table platforms
(
    id   bigserial not null primary key,
    name varchar(255) not null
);

COMMENT ON TABLE platforms IS 'platforms of candidates';
-- example comment with dataJSON
COMMENT ON COLUMN platforms.name IS 'full name';
-- examply index
create unique index platforms_name_idx
    ON platforms (name);

