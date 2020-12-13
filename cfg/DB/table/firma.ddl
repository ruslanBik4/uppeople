CREATE TABLE firma
(
    id                      serial,
    name                    character varying   not null,
    type_of_entity          character varying   not null DEFAULT '',
    the_company_fields      character varying   not null DEFAULT '',
    the_company_fields_memo text                not null DEFAULT '',
    description             character varying   not null DEFAULT '',
    bank                    text                NOT NULL DEFAULT '',
    addresses               character varying[] NOT NULL,
    emails                  character varying[] NOT NULL DEFAULT '{}',
    phones                  character varying[] NOT NULL DEFAULT '{}',
    memo                    text,
    EDPNOU                  integer,
    VAT                     integer,
    ITN                     bigint,
    IBAN                    character varying,
    PRIMARY KEY (id)
);
COMMENT ON TABLE firma IS 'Your Company information';
COMMENT ON COLUMN firma.name IS 'Name of the company';
COMMENT ON COLUMN firma.description IS 'Your Company information';
COMMENT ON COLUMN firma.type_of_entity IS 'Type of  entity';
COMMENT ON COLUMN firma.the_company_fields IS 'The company fields of activity';
COMMENT ON COLUMN firma.the_company_fields_memo IS 'Company fields of activity description';
COMMENT ON COLUMN firma.bank IS 'Bank account details';
COMMENT ON COLUMN firma.memo IS 'Note';
COMMENT ON COLUMN firma.EDPNOU IS 'code (EDPNOU) of legal entity {"pattern": "edpnou"}';
COMMENT ON COLUMN firma.VAT IS 'VAT payer registration, № {"pattern": "vat"}';
COMMENT ON COLUMN firma.ITN IS 'Tax payer number, ITN {"pattern": "itn"}';
COMMENT ON COLUMN firma.IBAN IS 'code (IBAN-) of bank deposit {"pattern": "iban"}';
COMMENT ON COLUMN firma.addresses IS 'The company addresses';
COMMENT ON COLUMN firma.emails IS 'list of emails company {"pattern": "email"}';
COMMENT ON COLUMN firma.phones IS 'list of phones company';
