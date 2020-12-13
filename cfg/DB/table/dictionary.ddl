CREATE TABLE dictionary
(
    id           serial,
    name         citext            not null,
    id_languages integer           not null,
    translation  character varying not null,
    PRIMARY KEY (id)
);
COMMENT ON TABLE dictionary IS 'word translate';
COMMENT ON COLUMN dictionary.name IS 'full name {"pattern": "notCyrillic"}';
COMMENT ON COLUMN dictionary.id_languages IS 'lang';
COMMENT ON COLUMN dictionary.translation IS 'name translation on language';
create unique index if not exists dictionary_name_idx
    ON dictionary (name, id_languages);
