CREATE OR REPLACE FUNCTION amoung_by_tags(sDate date, eDate date, userID integer, companyID integer,
            vacancyId integer, tags integer[])
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
       and create_at between COALESCE(sDate, NOW() - interval '1 month') and COALESCE( eDate, now() )
       and (companyID = 0 OR company_id = companyID)
       and (VacancyId = 0 OR vacancy_id = VacancyId)
       and (userID = 0 OR user_id = userID)
    into reContact;

    if vacancyId > 0 OR companyId > 0 then
        return query
            with rowsTags as (
                SELECT t.id::integer,
                       t.name,
                       t.color,
                       t.parent_id::integer,
                       (count(c.id) + CASE WHEN t.id::integer = 1 or t.name is null THEN reContact ELSE 0 END)::integer
                           as amount
                FROM tags t
                         JOIN candidates c ON t.id = c.tag_id
                         JOIN vacancies_to_candidates vtc on c.id = vtc.candidate_id
                         JOIN vacancies v ON v.id = vtc.vacancy_id
                WHERE c.date between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now())
                  and (companyID = 0 OR v.company_id = companyID)
                  and (vacancyId = 0 OR v.id = vacancyId
                    AND coalesce(vtc.date_last_change, vtc.date_create)
                                            between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now()))
                  and (userID = 0 OR c.recruter_id = userID)
                  and (tags is null or t.id = ANY (tags))
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
            ORDER BY 1 nulls last
        ;
    else
        return query
            with rowsTags as (
                SELECT t.id::integer,
                       t.name,
                       t.color,
                       t.parent_id::integer,
                       (count(c.id) + CASE WHEN t.id::integer = 1 or t.name is null THEN reContact ELSE 0 END)::integer
                           as amount
                FROM tags t
                         JOIN candidates c ON t.id = c.tag_id
                WHERE c.date between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now())
                  and (userID = 0 OR c.recruter_id = userID)
                  and (tags is null or t.id = ANY (tags))
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
            ORDER BY 1 nulls last
         ;
    END IF;
END;

$$;
