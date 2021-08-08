CREATE OR REPLACE PROCEDURE update_log()
LANGUAGE plpgsql
AS
$$
BEGIN
  update logs
  set changed=(
    CASE WHEN text LIKE '%=https:%' OR text LIKE '%comments=%' OR text LIKE '%details%' THEN null
         WHEN text LIKE '%=%' THEN
             (SELECT CONCAT('{',
                            array_to_string(
                                    array_agg(
                                            CONCAT('"',
                                                   REPLACE(
                                                           CASE
                                                               WHEN val LIKE '%[%]%'
                                                                   THEN REPLACE(val, ' ', ',')
                                                               WHEN val LIKE '%<p>%</p>%' THEN CONCAT(regexp_replace(
                                                                                                              REPLACE(REPLACE(REPLACE(val, '<p>', ''), '</p>', ''), '=', '="'),
                                                                                                              E'[\\n\\r\\f\\u000B\\u0085\\u2028\\u2029]+', '', 'g'), '"')
                                                               WHEN substring(val from '(?<==)[^\]]+') ~ '^0+\d+$' THEN CONCAT(REPLACE(val, '=', '="'), '"')
                                                               WHEN substring(val from '(?<==)[^\]]+') ~ '^\d+(\.\d)?$' THEN val
                                                               ELSE CONCAT(REPLACE(val, '=', '="'), '"') END
                                                       , '=', '":')
                                                )
                                        ), ','), '}')::jsonb
              FROM
                  regexp_split_to_table(text, E', ') as x(val))
         WHEN text LIKE '{%}' THEN text::jsonb
         ELSE to_jsonb(text) END)
where logs.text IS NOT NULL AND logs.changed IS NULL;
END;

$$;
