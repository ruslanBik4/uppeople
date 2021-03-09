CREATE TABLE users
(
    id           SERIAL              NOT NULL,
    name         character varying   NOT NULL,
    email        citext              not null,
    address      character varying   NOT NULL DEFAULT '',
    emails       character varying[] NOT NULL,
    phones       character varying[] NOT NULL,
    id_languages integer[],
    isDel        BOOLEAN             NOT NULL DEFAULT false,
    role_id     INTEGER             NOT NULL DEFAULT 3,
    last_login   timestamptz         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    hash         bigint              NOT NULL,
    last_page    character varying,
    id_homepages integer             NOT NULL DEFAULT 1,
    createAt     timestamp           not null default CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS plperlu;
COMMENT ON TABLE users IS 'Users profile data';
COMMENT ON COLUMN users.emails IS 'list of emails which notify user';
COMMENT ON COLUMN users.phones IS 'list of phones which contact user';
COMMENT ON COLUMN users.id_languages IS 'Preferred languages of user interface';
COMMENT ON COLUMN users.isDel IS 'Is deleting (read_only)';
COMMENT ON COLUMN users.hash IS 'Hash (read_only)';
COMMENT ON COLUMN users.role_id IS 'Role (read_only)';
COMMENT ON COLUMN users.id_homepages IS 'home page';
COMMENT ON COLUMN users.last_login IS 'Last time when user login (read_only)';
COMMENT ON COLUMN users.last_page IS 'Last page which visited user (read_only)';
COMMENT ON COLUMN users.createAt IS 'Date of user registration (read_only)';
COMMENT ON COLUMN users.email IS 'Email of user registration (read_only)';
create unique index if not exists users_name_idx
    ON users (email);
ALTER TABLE users
    ADD CONSTRAINT validemail_check CHECK
        (EMAIL ~
         '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$');
insert into users (name, email, hash, emails, phones)
values ('admin', 'zero@null.com', 0, '{}', '{}');
