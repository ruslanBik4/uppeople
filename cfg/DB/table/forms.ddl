CREATE TABLE forms
(
    id          serial,
    title       character varying not null,
    action      character varying not null,
    post        boolean           not null DEFAULT true,
    description character varying not null DEFAULT '',
    hideBlock   json,
    PRIMARY KEY (id)
);
COMMENT ON TABLE forms IS 'settings for forms generate';
-- example comment with dataJSON
COMMENT ON COLUMN forms.title IS 'full name {"pattern": "notCyrillic"}';
COMMENT ON COLUMN forms.hideBlock IS 'JSON of block hide/show {"pattern": "json"}';
-- examply index
create unique index if not exists forms_name_idx
    ON forms (title);

insert into public.forms (id, title, action, post, description)
values (1, 'form_editor', '/put', 'true', 'Edit from settings on  backend'),
       (3, 'fillers_adv', '/api/v1/table/fillers/put', 'true',
        'Exhausting information for all types of dispersed, reinforcing, and other types of fillers');
