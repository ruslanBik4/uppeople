CREATE TABLE form_fields
(
    id             serial,
    id_form_blocks integer           not null,
    name           character varying not null,
    input_type     character varying not null DEFAULT 'text',
    title          character varying not null,
    pattern        character varying not null DEFAULT '',
    placeholder    character varying not null DEFAULT '',
    autofocus      boolean           not null DEFAULT false,
    disabled       boolean           not null DEFAULT false,
    required       boolean           not null DEFAULT false,
    readOnly       boolean           not null DEFAULT false,
    PRIMARY KEY (id)
);
COMMENT ON TABLE form_fields IS 'fields of form';
-- example comment with dataJSON
COMMENT ON COLUMN form_fields.name IS 'name of parameter for sending form';
COMMENT ON COLUMN form_fields.title IS 'Label for input field';
-- examply index
create unique index if not exists form_fields_name_idx
    ON form_fields (id_form_blocks, name);


