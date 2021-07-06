CREATE TABLE log_actions
(
    id integer not null constraint log_actions_pkey primary key,
    name varchar not null,
    text_before_cand varchar not null,
    for_candidate varchar default '',
    for_company varchar default '',
    text_after_cand varchar default '',
    is_insert_text bool,
);

comment on table log_actions is 'actions logging derectives';

comment on column log_actions.name is 'const name';

create unique index log_actions_name_idx
    ON log_actions (name);

INSERT INTO log_actions (id,
                         name,
                         text_before_cand,
                         for_candidate,
                         for_company,
                         text_after_cand,
                         is_insert_text)
VALUES
(100, 'CODE_LOG_UPDATE', ' обновил/дополнил', ' у кандидата ', ' у компании ', ' следующую информацию:', true),
(101, 'CODE_LOG_INSERT', ' добавил', ' кандидата ', ' компанию ', '', false),
(102, 'CODE_LOG_PEFORM', ' совершил', ' по кандидату ', ' по компании ', ' следующее действие:', true),
(103, 'CODE_LOG_DELETE', ' удалил', ' кандидата ', ' контакт с компанией ', '', false),
(104, 'CODE_LOG_RE_CONTACT', ' обновил контакт', ' с кандидатом ', ' с компанией ', '', false),
(105, 'CODE_ADD_COMMENT', ' добавил/обновил комментарий', ' для кандидата ', ' для компании ', ':', true),
(106, 'CODE_DEL_COMMENT', ' удалил комментарий', ' для кандидата ', ' для компании ', ':', true),
(107, 'CODE_SEND_CV', ' отправил резюме', ' кандидата ', ' компании ', '', false),
(108, 'CODE_APPOINT_INTERVIEW', ' назначил интервью', ' кандидату ', ' компании ', '', false)