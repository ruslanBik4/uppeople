create table users
(
    id serial not null,
    role_id integer not null default 2,
    company_id integer,
    name character varying not null,
    email character varying not null,
    password character varying not null,
    phone character varying not null default ''::character varying,
    image character varying default NULL::character varying,
    user_freelancers integer,
    last_login timestamp not null default CURRENT_TIMESTAMP,
    last_ip character varying not null default '',
    PRIMARY KEY (id)
);

COMMENT ON TABLE users IS 'Users table with roles & other data';
COMMENT ON COLUMN users.name IS 'full name';
create unique index users_name_idx
    ON users (name, email);
create unique index users_email_idx
    ON users (email);
