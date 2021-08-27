CREATE OR REPLACE PROCEDURE update_log()
LANGUAGE plpgsql
AS
$$
BEGIN
  update logs
  set changed=(CASE WHEN text LIKE '%=https:%' OR text LIKE '%comments=%' OR text is null OR text = '' THEN null
         WHEN text LIKE '%=%' THEN
             (SELECT Format('{%s}',
                        string_agg(
                            Format(E'"%s',
                               REPLACE(
                                       CASE
                                           WHEN val = '' THEN 'none":""'
                                           WHEN val LIKE '%[%]%'
                                               THEN REPLACE(val, ' ', ',')
                                           WHEN val LIKE '%<p>%</p>%'
                                               THEN CONCAT(regexp_replace(
                                                          REPLACE(REPLACE(REPLACE(val, '<p>', ''), '</p>', ''), '=', '="'),
                                                          '/[\n\r\f\u000B\u0085\u2028\u2029]+/', '', 'g'), '"')
                                           WHEN substring(val from '(?<==)[^\]]+') ~ '^0+\d+$'
                                               THEN CONCAT(REPLACE(val, '=', '="'), '"')
                                           WHEN substring(val from '(?<==)[^\]]+') ~ '^\d+(\.\d)?$'
                                               THEN val
                                           ELSE CONCAT(REPLACE(val, '=', '="'), '"') END
                                   , '=', '":')
                                )
                        , ','))::jsonb
              FROM regexp_split_to_table(text, ',\s*') as x(val)
              where val > '')
         WHEN text LIKE '{%}' THEN text::jsonb
         ELSE to_jsonb(text)
        END)
--     from logs
    where logs.text is not null AND logs.changed IS NULL and vacancy_id is null;
END;

$$;
