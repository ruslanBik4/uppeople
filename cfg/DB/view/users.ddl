CREATE OR REPLACE VIEW lviv.users AS
SELECT *
FROM public.users
WHERE schema = 'lviv'
WITH CASCADED CHECK OPTION;

COMMENT ON VIEW lviv.users IS '';

