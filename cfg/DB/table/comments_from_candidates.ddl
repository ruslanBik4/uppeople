CREATE TABLE comments_from_candidates
(
    id           serial,
    user_id      integer default 1 not null,
    candidate_id integer           not null,
    comments     text,
    created_at   timestamp with time zone,
    updated_at   timestamp with time zone,
    PRIMARY KEY (id)
);

COMMENT ON TABLE comments_from_candidates IS '';

alter table comments_from_candidates
    add constraint comments_from_candidates_candidate_id_fk
        foreign key (candidate_id) references candidates
            on update cascade on delete set default;

-- alter table comments_from_candidates
--     add constraint comments_from_candidates_user_id_fk
--         foreign key (user_id) references users
--             on update cascade on delete set default;
