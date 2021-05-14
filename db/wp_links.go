// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_links struct {
	dbEngine.Table
	Record *Wp_linksFields
	rows   sql.Rows
}

type Wp_linksFields struct {
	Link_id          int64     `json:"link_id"`
	Link_url         string    `json:"link_url"`
	Link_name        string    `json:"link_name"`
	Link_image       string    `json:"link_image"`
	Link_target      string    `json:"link_target"`
	Link_description string    `json:"link_description"`
	Link_visible     string    `json:"link_visible"`
	Link_owner       float64   `json:"link_owner"`
	Link_rating      int64     `json:"link_rating"`
	Link_updated     time.Time `json:"link_updated"`
	Link_rel         string    `json:"link_rel"`
	Link_notes       string    `json:"link_notes"`
	Link_rss         string    `json:"link_rss"`
}

func (r *Wp_linksFields) RefColValue(name string) interface{} {
	switch name {
	case "link_id":
		return &r.Link_id

	case "link_url":
		return &r.Link_url

	case "link_name":
		return &r.Link_name

	case "link_image":
		return &r.Link_image

	case "link_target":
		return &r.Link_target

	case "link_description":
		return &r.Link_description

	case "link_visible":
		return &r.Link_visible

	case "link_owner":
		return &r.Link_owner

	case "link_rating":
		return &r.Link_rating

	case "link_updated":
		return &r.Link_updated

	case "link_rel":
		return &r.Link_rel

	case "link_notes":
		return &r.Link_notes

	case "link_rss":
		return &r.Link_rss

	default:
		return nil
	}
}

func (r *Wp_linksFields) ColValue(name string) interface{} {
	switch name {
	case "link_id":
		return r.Link_id

	case "link_url":
		return r.Link_url

	case "link_name":
		return r.Link_name

	case "link_image":
		return r.Link_image

	case "link_target":
		return r.Link_target

	case "link_description":
		return r.Link_description

	case "link_visible":
		return r.Link_visible

	case "link_owner":
		return r.Link_owner

	case "link_rating":
		return r.Link_rating

	case "link_updated":
		return r.Link_updated

	case "link_rel":
		return r.Link_rel

	case "link_notes":
		return r.Link_notes

	case "link_rss":
		return r.Link_rss

	default:
		return nil
	}
}

func NewWp_links(db *dbEngine.DB) (*Wp_links, error) {
	table, ok := db.Tables[TABLE_WPLinks]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_WPLinks}
	}

	return &Wp_links{
		Table: table,
	}, nil
}

func (t *Wp_links) NewRecord() *Wp_linksFields {
	t.Record = &Wp_linksFields{}
	return t.Record
}

func (t *Wp_links) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_links) SelectSelfScanEach(ctx context.Context, each func(record *Wp_linksFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_links) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_links) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
