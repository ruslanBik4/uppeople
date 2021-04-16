CREATE OR REPLACE FUNCTION amoung_by_status(sDate date, eDate date, userID integer, companyID integer,
            vacancyId integer, statuses integer[])
    RETURNS table(
                    id integer,
                    name character varying,
                    color character varying,
                    count integer,
                    percent numeric(5,2)
                 )
    LANGUAGE plpgsql
AS
$$
BEGIN
    return query
      with rowsStatus as (
          SELECT vtc.status, sfv.status as name, sfv.color, count(vtc.id)::integer as amount
          FROM vacancies_to_candidates vtc
                   JOIN candidates c ON c.id = vtc.candidate_id
                   JOIN status_for_vacs sfv ON sfv.id = vtc.status
                   JOIN vacancies v ON v.id = vtc.vacancy_id
          WHERE c.date between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now())
            AND coalesce(vtc.date_last_change, vtc.date_create)
              between COALESCE(sDate, NOW() - interval '1 month') and COALESCE(eDate, now())
            and (companyID = 0 OR v.company_id = companyID)
            and (vacancyId = 0 OR v.id = vacancyId)
            and (userID = 0 OR c.recruter_id = userID)
            and (statuses is null OR vtc.status = ANY (statuses))
          GROUP BY grouping sets ((1,2,3),())
          ORDER BY 1 nulls last
      )
      select *, ((amount * 100)::numeric / (select amount from rowsStatus where status is null))::numeric(5,2)
      from rowsStatus;
END;
$$;
