CREATE OR REPLACE VIEW select_vacancies AS
select v.id,
       label,
       LOWER(label) as value,
       v.platform_id,
       v.user_ids, v.status
FROM vacancies v,
     LATERAL format('%s (%L) %s',

                    (select c.name from companies c where v.company_id= c.id),
                    (select p.name from public.platforms p where v.platform_id = p.id),
                    (select s.name from public.seniorities s where s.id = v.seniority_id)
            ) as label;

COMMENT ON VIEW select_vacancies IS 'vacancies object for select view/edit on forms';

