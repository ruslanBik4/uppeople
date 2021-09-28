CREATE OR REPLACE VIEW lviv.users AS
SELECT *
FROM public.users
WHERE schema = 'lviv'
WITH CASCADED CHECK OPTION;

COMMENT ON VIEW lviv.users IS '';

CREATE OR REPLACE FUNCTION lviv.insert_user()
    RETURNS trigger
    LANGUAGE plpgsql
AS
$$
declare
    columns text;
    sql     text;
BEGIN
    NEW.schema = 'lviv';

    INSERT INTO public.users (name, email, phone, role_id, schema, hash)
            VALUES (NEW.name, NEW.email, NEW.phone, NEW.role_id, NEW.schema, NEW.hash)
            RETURNING id INTO NEW.id;
    return NEW;
END;
$$;

create trigger insert_users_trg
    instead of insert
    on lviv.users
    for each row
execute function insert_user();
