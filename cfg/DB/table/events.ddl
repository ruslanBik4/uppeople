CREATE TABLE events
(
    id          serial,
    id_firma    integer           not null,
    name        character varying not null,
    description character varying not null DEFAULT '',
    memo        text,
    place       character varying not null DEFAULT '',
    start_date  date              not null,
    end_date    date              not null,
    id_photos   integer           not null,
    PRIMARY KEY (id)
);

comment on table events is 'Form for the events';
COMMENT ON COLUMN events.name IS 'The name of the event organizing company';
COMMENT ON COLUMN events.description IS 'Information about the event';
COMMENT ON COLUMN events.place IS 'The event address';
COMMENT ON COLUMN events.start_date IS 'The event start date';
COMMENT ON COLUMN events.end_date IS 'The event end date';
