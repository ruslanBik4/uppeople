CREATE TABLE patterns_list
(
    id          serial            not null
        constraint patterns_list_pkey primary key,
    name        character varying not null,
    pattern     character varying not null,
    description character varying not null DEFAULT ''
);
comment on table patterns_list is 'List of pattern which help us to show edit of fields';
COMMENT ON COLUMN patterns_list.name IS 'Code name {"pattern": "notCyrillic"}';
COMMENT ON COLUMN patterns_list.pattern IS 'regular express {"pattern": "^[^\s]*$"}';
COMMENT ON COLUMN patterns_list.description IS 'User friendly explanation about input according to pattern';
create unique index if not exists patterns_list_name_idx
    on patterns_list (name);

INSERT INTO patterns_list (name, pattern, description)
VALUES ('inn', '\d{10,12}', ''),
       ('name', '^[^\s]+$', ''),
       ('notCyrillic', '^[\x1F-\xBF]*$', 'only latin word + digitals'),
       ('floatPos', '^\d+(\.\d{1,2})?$', 'only float number more than zero'),
       ('index', '[0-9]{6}', ''),
       ('ip', '\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', ''),
       ('phone', '(^\+[0-9]{1,3}\s?\([0-9]{1,2}\)\s?[0-9-]{1,9})$', ''),
       ('bankCard', '\d{16}', 'number of bank card');
