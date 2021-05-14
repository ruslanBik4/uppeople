// generate file
// don't edit
package db

import (
	"database/sql"
	"time"
	"unsafe"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_simply_static_pages struct {
	dbEngine.Table
	Record *Wp_simply_static_pagesFields
	rows   sql.Rows
}

type Wp_simply_static_pagesFields struct {
	Id                  int64           `json:"id"`
	Found_on_id         sql.NullFloat64 `json:"found_on_id"`
	Url                 string          `json:"url"`
	Redirect_url        sql.NullString  `json:"redirect_url"`
	File_path           sql.NullString  `json:"file_path"`
	Http_status_code    sql.NullInt32   `json:"http_status_code"`
	Content_type        sql.NullString  `json:"content_type"`
	Content_hash        unsafe.Pointer  `json:"content_hash"`
	Error_message       sql.NullString  `json:"error_message"`
	Status_message      sql.NullString  `json:"status_message"`
	Last_checked_at     time.Time       `json:"last_checked_at"`
	Last_modified_at    time.Time       `json:"last_modified_at"`
	Last_transferred_at time.Time       `json:"last_transferred_at"`
	Created_at          time.Time       `json:"created_at"`
	Updated_at          time.Time       `json:"updated_at"`
}

func (r *Wp_simply_static_pagesFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "found_on_id":
		return &r.Found_on_id

	case "url":
		return &r.Url

	case "redirect_url":
		return &r.Redirect_url

	case "file_path":
		return &r.File_path

	case "http_status_code":
		return &r.Http_status_code

	case "content_type":
		return &r.Content_type

	case "content_hash":
		return &r.Content_hash

	case "error_message":
		return &r.Error_message

	case "status_message":
		return &r.Status_message

	case "last_checked_at":
		return &r.Last_checked_at

	case "last_modified_at":
		return &r.Last_modified_at

	case "last_transferred_at":
		return &r.Last_transferred_at

	case "created_at":
		return &r.Created_at

	case "updated_at":
		return &r.Updated_at

	default:
		return nil
	}
}

func (r *Wp_simply_static_pagesFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "found_on_id":
		return r.Found_on_id

	case "url":
		return r.Url

	case "redirect_url":
		return r.Redirect_url

	case "file_path":
		return r.File_path

	case "http_status_code":
		return r.Http_status_code

	case "content_type":
		return r.Content_type

	case "content_hash":
		return r.Content_hash

	case "error_message":
		return r.Error_message

	case "status_message":
		return r.Status_message

	case "last_checked_at":
		return r.Last_checked_at

	case "last_modified_at":
		return r.Last_modified_at

	case "last_transferred_at":
		return r.Last_transferred_at

	case "created_at":
		return r.Created_at

	case "updated_at":
		return r.Updated_at

	default:
		return nil
	}
}

func NewWp_simply_static_pages(db *dbEngine.DB) (*Wp_simply_static_pages, error) {
	table, ok := db.Tables[TABLE_WPSimplyStaticPages]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_WPSimplyStaticPages}
	}

	return &Wp_simply_static_pages{
		Table: table,
	}, nil
}

func (t *Wp_simply_static_pages) NewRecord() *Wp_simply_static_pagesFields {
	t.Record = &Wp_simply_static_pagesFields{}
	return t.Record
}

func (t *Wp_simply_static_pages) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_simply_static_pages) SelectSelfScanEach(ctx context.Context, each func(record *Wp_simply_static_pagesFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_simply_static_pages) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_simply_static_pages) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
