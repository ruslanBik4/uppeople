create table tags
(
    id serial not null,
    name character varying not null,
    color character varying not null,
    parent_id integer not null,
    order_num integer not null,
    PRIMARY KEY (id)
);

COMMENT ON TABLE tags IS 'tags of candidates state';
-- example comment with dataJSON
COMMENT ON COLUMN tags.name IS 'full name';
-- examply index
create unique index tags_name_idx
    ON tags (name);
