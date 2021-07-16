CREATE OR REPLACE FUNCTION json_array_castint(json) RETURNS int[] AS $f$
SELECT (CASE WHEN $1::text NOT LIKE '[%]' THEN (SELECT array_agg($1::text::int)::int[])
    ELSE
(SELECT array_agg(x)::int[] || ARRAY[]::int[] FROM json_array_elements_text($1) t(x)) END);
$f$ LANGUAGE sql IMMUTABLE;