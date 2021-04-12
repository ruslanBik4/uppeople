CREATE TABLE form_blocks
(
    id          serial,
    id_forms    integer           not null,
    title       character varying not null,
    description character varying not null DEFAULT '',
    tableName   character varying not null,
    columns     character varying[],
    buttons     JSON,
    PRIMARY KEY (id)
);
COMMENT ON TABLE form_blocks IS 'part of form setting';
COMMENT ON COLUMN form_blocks.title IS 'Title of fields block {"pattern": "notCyrillic"}';
COMMENT ON COLUMN form_blocks.description IS 'Description of fields block';
COMMENT ON COLUMN form_blocks.tableName IS 'Name of table that must populate from block fields';
COMMENT ON COLUMN form_blocks.columns IS 'list of tables columns which generate inputs field of block';
COMMENT ON COLUMN form_blocks.buttons IS 'JSON of blocks {"pattern": "json"}';
create unique index if not exists form_blocks_name_idx
    ON form_blocks (id_forms, title);

insert into form_blocks (id, title, description, buttons, id_forms, tablename)
values (1, 'Form setting', 'form properties', '{
  "title": "Send fs"
}', 1, 'forms'),
       (2, 'Block setting', 'block properties', '{
         "title": "Send fs"
       }', 1, 'form_blocks'),
       (3, 'Fields setting', 'fields properties', '{
         "title": "Send fs"
       }', 1, 'form_fields');
insert into form_blocks (id, title, description, buttons, id_forms, tablename, columns)
values (4, 'Information', '', null, 3, 'fillers', '{id_filler_for_polymers, marka,
       manufacturer,
       volume_on_stock}'),
       (5, 'Показник властивостей', 'та стандарти випробувань', null, 3, 'fillers',
        '{test_method,
              chemical_nature,
              content_filler,
              density,
              filament_diameter,
              number_of_filaments,
              id_type_of_filler}'),
       (6, 'Діаметр елементарних ниток, мкм', '', null, 3, 'fillers', '{length_of_fibres}'),
       (7, 'Дисперсність, мкм', '', null, 3, 'fillers', '{dispersity}'),
       (8, '', '', null, 3, 'fillers', '{linear_density,
       moisture_content,
       purity_of_disperse_filler,
       tensile_strength,
       tensile_modulus,
       strain_failure,
       tds,
       msds,
       presentation}')
;

