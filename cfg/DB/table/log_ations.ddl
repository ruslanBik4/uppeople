CREATE TABLE log_actions
(
    id integer not null constraint log_actions_pkey primary key,
    name varchar not null,
    text_before_cand varchar not null,
    text_after_cand varchar default '',
    --format_str varchar default '',
    is_insert_text bool,
    for_candidate varchar default '',
    for_company varchar default '',
);

comment on table log_actions is 'actions logging derectives';

comment on column log_actions.name is 'const name';

create unique index log_actions_name_idx
    ON log_actions (name);