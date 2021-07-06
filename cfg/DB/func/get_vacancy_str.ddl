CREATE OR REPLACE FUNCTION get_vacancy_str(Id integer)
    RETURNS table (text text)
    LANGUAGE plpgsql
AS
$$
BEGIN
return query
select CONCAT(platforms.name, ', ', seniorities.name, ' (',
              CASE WHEN vacancies.name is not null
                       THEN vacancies.name
                   ELSE '' END,
              companies.name, ')')

from vacancies
         left Join companies on (vacancies.company_id = companies.id)
         left Join platforms on (vacancies.platform_id = platforms.id)
         left Join seniorities on (vacancies.seniority_id = seniorities.id)

where vacancies.id = $1
;
END;
$$;