CREATE OR REPLACE FUNCTION amoung_by_tags(sDate date, eDate date, userID integer, companyID integer, VacancyId integer)
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
DECLARE
BEGIN

    if vacancyId > 0 OR companyId > 0 then
        return query
        SELECT t.id, t.name, t.color,
               count(c.id) + (select count(*)
                              from logs
                              where kod_deystviya = 104
                               and create_at between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
                                and (companyID == 0 OR company_id = companyID)
                                and (VacancyId == 0 OR vacancy_id = VacancyId)
                                and (userID == 0 OR user_id = userID)
               ) as count,
               t.parent_id
        FROM tags t JOIN candidates c ON t.id=c.tag_id
                    JOIN vacancies_to_candidates vtc on c.id = vtc.candidate_id
                    JOIN vacancies v ON v.id = vtc.vacancy_id
        WHERE c.date  between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
          and (companyID == 0 OR v.company_id = companyID)
          and (VacancyId == 0 OR v.id = VacancyId)
          and (userID == 0 OR c.recruter_id = userID)
        GROUP BY 1, 2, 3, 5
        ORDER BY 1;
    else
        return query
        SELECT t.id, t.name, t.color,
               count(c.id) + (select count(*)
                              from logs
                              where kod_deystviya = 104
                               and create_at between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now())
               ) as count,
               t.parent_id
        FROM tags t JOIN candidates c ON t.id=c.tag_id
        WHERE c.date  between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
          and (userID == 0 OR c.recruter_id = userID)
        GROUP BY 1, 2, 3, 5
        ORDER BY 1;
    END IF;
END;

$$;
