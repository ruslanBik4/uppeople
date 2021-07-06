CREATE TABLE patterns_list
(
    id          serial,
    name        character varying not null,
    pattern     character varying not null,
    description character varying not null DEFAULT '',
    PRIMARY KEY (id)
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
       ('notCyrillic', '^[\x1F-\xBF]*$', 'only latin word + digital'),
       ('floatPos', '^\d+(\.\d{1,2})?$', 'only float number more than zero'),
       ('index', '[0-9]{6}', ''),
       ('ip', '\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}', ''),
       ('email',
        '^([^<>()[\\]\\.,;:\\s@]+(\\.[^<>()[\\]\\.,;:\\s@]+)*@(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))$',
        'Correct email need'
       ),
       ('phone', '(^\+[0-9]{1,3}\s?\([0-9]{1,2}\)\s?[0-9-]{1,9})$', ''),
       ('bankCard', '\d{16}', 'number of bank card');
--        ('full_email',
--         '^(([^<>()[\]\.,;:\s@""]+(\.[^<>()[\]\.,;:\s@""]+)*)|("".+""))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$',
--         'Must correct email, ex. ruslan@pm-db.net'),
-- +380(77)89898332