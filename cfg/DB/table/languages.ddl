CREATE TABLE languages
(
    id   serial,
    name char(2) not null,
    PRIMARY KEY (id)
);
COMMENT ON TABLE languages IS '';
COMMENT ON COLUMN languages.name IS 'short name of language';
create unique index if not exists languages_name_idx
    ON languages (name);
