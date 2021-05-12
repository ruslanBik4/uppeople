CREATE OR REPLACE FUNCTION amoung_by_tags(sDate date, eDate date, userID integer, companyID integer,
                                          platformId integer, vacancyId integer, tags integer[])
    RETURNS table(
                    id integer,
                    name character varying,
                    color character varying,
                    parent_id integer,
                    count integer,
                    percent numeric(5,2)
                 )
    LANGUAGE plpgsql
AS
$$
DECLARE reContact integer = 0;
BEGIN
    select count(*)
     from logs
     where kod_deystviya = 104
       and create_at ::date between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
       and (companyID = 0 OR company_id = companyID)
       and (vacancyId = 0 OR vacancy_id = vacancyId)
       and (platformId = 0 OR exists(select NULL from vacancies v where platform_id = platformId AND v.id=vacancy_id) )
       and (userID = 0 OR user_id = userID)
    into reContact;

    if vacancyId > 0 OR companyId > 0 OR platformId > 0 then
        return query
            with rowsTags as (
                SELECT t.id,
                       t.name,
                       t.color,
                       t.parent_id,
                       (count(c.id) + CASE WHEN t.id = 1 or t.name is null THEN reContact ELSE 0 END)::integer
                           as amount
                FROM tags t
                         JOIN candidates c ON t.id = c.tag_id
                         LEFT JOIN vacancies_to_candidates vtc on c.id = vtc.candidate_id
                         JOIN vacancies v ON (v.id = vtc.vacancy_id
                                                OR (vtc.vacancy_id is null AND v.id = ANY (c.vacancies))
                                            )
                WHERE c.date ::date between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now())
                  and (companyID = 0 OR v.company_id = companyID)
                  and (vacancyId = 0 OR v.id = vacancyId
                    AND coalesce(vtc.date_last_change, vtc.date_create, c.date) ::date
                                            between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now()))
                  and (platformId = 0 OR v.platform_id = platformId)
                  and (userID = 0 OR c.recruter_id = userID)
                  and (tags is null
                    OR (t.parent_id = 0 AND t.id = ANY (tags) OR t.parent_id=ANY(tags)))
                GROUP BY grouping sets ((1, 2, 3, 4), ())
            )
            select t.id,
                   t.name,
                   t.color,
                   t.parent_id,
                   sum(amount)::integer,
                   ((sum(amount) * 100)::numeric(8,2) / (select amount from rowsTags where rowsTags.id is null))::numeric(5,2)
            from rowsTags r JOIN tags t ON t.id = r.parent_id
            where r.parent_id > 0
            GROUP BY 1, 2, 3, 4
            union
            select *,
                   ((amount * 100)::numeric(8,2) / (select amount from rowsTags where rowsTags.id is null))::numeric(5,2)
            from rowsTags
            where amount > 0
            ORDER BY 1 nulls last
        ;
    else
        return query
            with rowsTags as (
                SELECT t.id,
                       t.name,
                       t.color,
                       t.parent_id,
                       (count(c.id) + CASE WHEN t.id = 1 or t.name is null THEN reContact ELSE 0 END)::integer
                           as amount
                FROM tags t
                         JOIN candidates c ON t.id = c.tag_id
                WHERE c.date ::date between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now())
                  and (userID = 0 OR c.recruter_id = userID)
                  and (tags is null
                       OR (t.parent_id = 0 AND t.id = ANY (tags) OR t.parent_id=ANY(tags)))
                GROUP BY grouping sets ((1, 2, 3, 4), ())
            )
            select t.id,
                   t.name,
                   t.color,
                   t.parent_id,
                   sum(amount)::integer,
                   ((sum(amount) * 100)::numeric(8,2) / (select amount from rowsTags where rowsTags.id is null))::numeric(5,2)
            from rowsTags r JOIN tags t ON t.id = r.parent_id
            where r.parent_id > 0
            GROUP BY 1, 2, 3, 4
            union
            select *,
                    ((amount * 100)::numeric(8,2) / (select amount from rowsTags where rowsTags.id is null))::numeric(5,2)
            from rowsTags
            where amount > 0
            ORDER BY 1 nulls last
         ;
    END IF;
END;

$$;
