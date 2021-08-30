CREATE OR REPLACE VIEW all_staff AS
SELECT *
FROM public.users u,
     LATERAL (select count(*) FILTER ( WHERE action_code = (select id from public.log_actions where name = 'CODE_LOG_INSERT') ) create_count,
                     count(*) FILTER ( WHERE action_code = (select id from public.log_actions where name = 'CODE_LOG_UPDATE') ) update_count,
                     count(*) FILTER ( WHERE action_code = (select id from public.log_actions where name = 'CODE_LOG_PEFORM') ) send_count
              from logs
              where age(date_create) < interval '7 day' and user_id = u.id ) j;

COMMENT ON VIEW all_staff IS 'view all users with statistics';

