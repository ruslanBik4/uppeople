CREATE TABLE homepages
(
    id   serial,
    name character varying not null,
    url  character varying not null,
    PRIMARY KEY (id)
);
COMMENT ON TABLE homepages IS 'list of available URL homepages';
-- example comment with dataJSON
COMMENT ON COLUMN homepages.name IS 'full name {"pattern": "name"}';
COMMENT ON COLUMN homepages.url IS 'full name {"pattern": "url"}';
-- examply index
create unique index if not exists homepages_name_idx
    ON homepages (name);
insert into homepages (name, url)
values ('profile', '/user/profile')