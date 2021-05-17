CREATE TABLE status_for_vacs
(
    id   serial,
    status character varying not null,
    color character varying not null,
    order_num integer not null,
    PRIMARY KEY (id)
);
COMMENT ON TABLE status_for_vacs IS 'Statuses of candidates on vacancies';
create unique index status_for_vacs_name_idx
    ON status_for_vacs (status);
