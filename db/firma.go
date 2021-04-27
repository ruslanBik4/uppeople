// generate file
// don't edit
package db

import (
	"database/sql"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Firma struct {
	dbEngine.Table
	Record *FirmaFields
	rows   sql.Rows
}

type FirmaFields struct {
	Id                      int32          `json:"id"`
	Name                    string         `json:"name"`
	Type_of_entity          string         `json:"type_of_entity"`
	The_company_fields      string         `json:"the_company_fields"`
	The_company_fields_memo string         `json:"the_company_fields_memo"`
	Description             string         `json:"description"`
	Bank                    string         `json:"bank"`
	Addresses               []string       `json:"addresses"`
	Emails                  []string       `json:"emails"`
	Phones                  []string       `json:"phones"`
	Memo                    sql.NullString `json:"memo"`
	Edpnou                  sql.NullInt32  `json:"edpnou"`
	Vat                     sql.NullInt32  `json:"vat"`
	Itn                     sql.NullInt64  `json:"itn"`
	Iban                    sql.NullString `json:"iban"`
}

func (r *FirmaFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "name":
		return &r.Name

	case "type_of_entity":
		return &r.Type_of_entity

	case "the_company_fields":
		return &r.The_company_fields

	case "the_company_fields_memo":
		return &r.The_company_fields_memo

	case "description":
		return &r.Description

	case "bank":
		return &r.Bank

	case "addresses":
		return &r.Addresses

	case "emails":
		return &r.Emails

	case "phones":
		return &r.Phones

	case "memo":
		return &r.Memo

	case "edpnou":
		return &r.Edpnou

	case "vat":
		return &r.Vat

	case "itn":
		return &r.Itn

	case "iban":
		return &r.Iban

	default:
		return nil
	}
}

func (r *FirmaFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "name":
		return r.Name

	case "type_of_entity":
		return r.Type_of_entity

	case "the_company_fields":
		return r.The_company_fields

	case "the_company_fields_memo":
		return r.The_company_fields_memo

	case "description":
		return r.Description

	case "bank":
		return r.Bank

	case "addresses":
		return r.Addresses

	case "emails":
		return r.Emails

	case "phones":
		return r.Phones

	case "memo":
		return r.Memo

	case "edpnou":
		return r.Edpnou

	case "vat":
		return r.Vat

	case "itn":
		return r.Itn

	case "iban":
		return r.Iban

	default:
		return nil
	}
}

func NewFirma(db *dbEngine.DB) (*Firma, error) {
	table, ok := db.Tables[TableFirma]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TableFirma}
	}

	return &Firma{
		Table: table,
	}, nil
}

func (t *Firma) NewRecord() *FirmaFields {
	t.Record = &FirmaFields{}
	return t.Record
}

func (t *Firma) GetFields(columns []dbEngine.Column) []interface{} {
	if len(columns) == 0 {
		columns = t.Columns()
	}

	t.NewRecord()
	v := make([]interface{}, len(columns))
	for i, col := range columns {
		v[i] = t.Record.RefColValue(col.Name())
	}

	return v
}

func (t *Firma) SelectSelfScanEach(ctx context.Context, each func(record *FirmaFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Firma) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
	if len(Options) == 0 {
		v := make([]interface{}, len(t.Columns()))
		columns := make([]string, len(t.Columns()))
		for i, col := range t.Columns() {
			columns[i] = col.Name()
			v[i] = t.Record.ColValue(col.Name())
		}
		Options = append(Options,
			dbEngine.ColumnsForSelect(columns...),
			dbEngine.ArgsForSelect(v...))
	}

	return t.Table.Insert(ctx, Options...)
}

func (t *Firma) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
	if len(Options) == 0 {
		v := make([]interface{}, len(t.Columns()))
		priV := make([]interface{}, 0)
		columns := make([]string, 0, len(t.Columns()))
		priColumns := make([]string, 0, len(t.Columns()))
		for _, col := range t.Columns() {
			if col.Primary() {
				priColumns = append(priColumns, col.Name())
				priV[len(priColumns)-1] = t.Record.ColValue(col.Name())
				continue
			}

			columns = append(columns, col.Name())
			v[len(columns)-1] = t.Record.ColValue(col.Name())
		}

		Options = append(
			Options,
			dbEngine.ColumnsForSelect(columns...),
			dbEngine.WhereForSelect(priColumns...),
			dbEngine.ArgsForSelect(append(v, priV...)...),
		)
	}

	return t.Table.Update(ctx, Options...)
}
