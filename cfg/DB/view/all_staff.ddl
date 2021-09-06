CREATE OR REPLACE VIEW lviv.all_staff AS
SELECT *
FROM public.users u,
     LATERAL (select count(*) FILTER ( WHERE name = 'CODE_LOG_INSERT') create_count,
                     count(*) FILTER ( WHERE name = 'CODE_LOG_UPDATE') update_count,
                     count(*) FILTER ( WHERE name = 'CODE_LOG_PEFORM') send_count
              from logs join public.log_actions a on (action_code = a.id)
              where age(date_create) < interval '7 day' and user_id = u.id ) j
WHERE schema = 'lviv';

COMMENT ON VIEW all_staff IS 'view all users with statistics';

