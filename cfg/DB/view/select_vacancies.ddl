CREATE OR REPLACE VIEW select_vacancies AS
select v.id,
       label,
       LOWER(label) as value,
       v.platform_id,
       v.user_ids
FROM vacancies v JOIN companies c on (v.company_id= c.id)
                 JOIN platforms p ON (v.platform_id = p.id),
                 LATERAL concat(c.name, ' ("', p.name, '")') as label;

COMMENT ON VIEW select_vacancies IS 'vacancies object for select view/edit on forms';

