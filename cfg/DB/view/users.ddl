CREATE OR REPLACE VIEW lviv.users AS
SELECT *
FROM public.users
WHERE schema = 'lviv';
COMMENT ON VIEW users IS '';
COMMENT ON COLUMN users.* IS 'full name {"pattern": "name"}';

