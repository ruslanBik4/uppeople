CREATE OR REPLACE FUNCTION whole_report(sDate date, eDate date)
    RETURNS table(
                     id integer,
                     name character varying,
                     company character varying,
                     vacancy character varying,
                     recruter character varying,
                     count integer,
                     percent numeric(5,2)
                 )
        LANGUAGE plpgsql
AS
$$
BEGIN

    return query
    with v as (
        SELECT t.id,
               t.name,
               v.company_id,
               c.recruter_id,
               v.id as vacancy_id,
               CONCAT( (select name from platforms p where p.id = v.platform_id),
                       ' (',
                       (select s.status from statuses s where s.id = v.status),
                       ')'
                   ) as vacancy,
               count(c.id)::integer  as amount
        FROM tags t
                 JOIN candidates c ON t.id = c.tag_id
                 LEFT JOIN vacancies_to_candidates vtc on c.id = vtc.candidate_id
                 JOIN vacancies v ON (v.id = vtc.vacancy_id
            OR (vtc.vacancy_id is null AND v.id = ANY (C.vacancies))
            )
        WHERE c.date between COALESCE($1, NOW() - interval '1 month') and COALESCE($2, now())
          AND coalesce(vtc.date_last_change, vtc.date_create, c.date)
            between COALESCE($1, NOW() - interval '1 month') and COALESCE($2, now())
        GROUP BY CUBE  ((1,2), 3, 4, (5, 6))
    )
    select v.id, v.name,
           coalesce((select co.name from companies co where co.id= v.company_id), 'total') as company,
           v.vacancy :: character varying,
           (select u.name from users u where u.id = recruter_id) as recruter,
           v.amount,
           CASE WHEN v.id = 1 THEN
                    (select count(*)
                     from logs
                     where kod_deystviya = 104
                       and create_at between COALESCE($1, NOW() - interval '1 month') and COALESCE( $2, now() )
                       and (v.company_id is null OR company_id = v.company_id)
                       and (v.vacancy_id is null OR vacancy_id = v.vacancy_id)
                       and (v.recruter_id is null OR user_id = v.recruter_id))
               ELSE 0
               END
               :: numeric(8,2)                                                                    as reContact
    from v
    ORDER BY 1,2,3,4,5 nulls last;

END;

$$;
