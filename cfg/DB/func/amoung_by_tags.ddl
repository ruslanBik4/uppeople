CREATE OR REPLACE FUNCTION amoung_by_tags(sDate date, eDate date, userID integer, companyID integer, vacancyId integer)
    RETURNS table(
                    id integer,
                    name character varying,
                    color character varying,
                    count integer,
                    parent_id integer
                 )
    LANGUAGE plpgsql
AS
$$
DECLARE reContact integer = 0;
BEGIN
    select count(*)
     from logs
     where kod_deystviya = 104
       and create_at between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
       and (companyID = 0 OR company_id = companyID)
       and (VacancyId = 0 OR vacancy_id = VacancyId)
       and (userID = 0 OR user_id = userID)
    into reContact;

    if vacancyId > 0 OR companyId > 0 then
        return query
        SELECT t.id::integer, t.name, t.color,
               (count(c.id) + CASE WHEN t.id::integer = 1 THEN reContact ELSE 0 END)::integer as count,
               t.parent_id::integer
        FROM tags t JOIN candidates c ON t.id=c.tag_id
                        OR c.tag_id in (select t2.id from tags t2 where t2.parent_id = t.id)
                    JOIN vacancies_to_candidates vtc on c.id = vtc.candidate_id
                    JOIN vacancies v ON v.id = vtc.vacancy_id
        WHERE c.date between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
          and (companyID = 0 OR v.company_id = companyID)
          and (vacancyId = 0 OR v.id = vacancyId
              AND coalesce(vtc.date_last_change, vtc.date_create)
                  between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() ))
          and (userID = 0 OR c.recruter_id = userID)
        GROUP BY 1, 2, 3, 5
        ORDER BY 1;
    else
        return query
        SELECT t.id::integer, t.name, t.color,
               (count(c.id) + CASE WHEN t.id::integer = 1 THEN reContact ELSE 0 END)::integer as count,
               t.parent_id::integer
        FROM tags t
            JOIN candidates c ON t.id=c.tag_id
                            OR c.tag_id in (select t2.id from tags t2 where t2.parent_id = t.id)
        WHERE c.date  between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
          and (userID = 0 OR c.recruter_id = userID)
        GROUP BY 1, 2, 3, 5
        ORDER BY 1;
    END IF;
END;

$$;
