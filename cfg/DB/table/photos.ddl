CREATE TABLE photos
(
    id   serial            not null
        constraint photos_pkey primary key,
    name character varying not null,
    blob bytea             not null
);
comment on table photos is 'Store of all photo';
CREATE EXTENSION IF NOT EXISTS pgcrypto;

create unique index if not exists photos_blob_uindex
    on photos (digest(blob, 'sha1'));

