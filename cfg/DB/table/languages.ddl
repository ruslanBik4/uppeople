CREATE TABLE languages
(
    id   serial,
    name character varying not null,
    abbr char(2) not null,
    PRIMARY KEY (id)
);
COMMENT ON TABLE languages IS '';
COMMENT ON COLUMN languages.name IS 'name of language';
COMMENT ON COLUMN languages.abbr IS 'short name of language';
create unique index languages_name_idx
    ON languages (name);
