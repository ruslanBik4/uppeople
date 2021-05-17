create table seniorities
(
    id    serial not null,
    name varchar(255) not null,
    primary key (id)
);

COMMENT ON TABLE seniorities IS 'senoirities of candidates';
-- example comment with dataJSON
COMMENT ON COLUMN seniorities.name IS 'full name';
-- examply index
create unique index seniorities_name_idx
    ON seniorities (name);