CREATE OR REPLACE FUNCTION jsonb_array_castint(jsonb) RETURNS int[] AS $f$
SELECT (CASE WHEN $1::text NOT LIKE '[%]' THEN (SELECT array_agg($1::text::int)::int[])
        ELSE
        (SELECT array_agg(x)::int[] || ARRAY[]::int[] FROM jsonb_array_elements_text($1) t(x)) END);
$f$ LANGUAGE sql IMMUTABLE;