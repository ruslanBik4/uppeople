CREATE TABLE management
(
    id_users integer not null,
    id_firma integer not null,
    grants   JSON,
    PRIMARY KEY (id_users, id_firma)
);
COMMENT ON TABLE management IS 'List of grands';
COMMENT ON COLUMN management.id_users IS 'Users id';
COMMENT ON COLUMN management.id_firma IS 'Company';
COMMENT ON COLUMN management.grants IS 'Grands list';

