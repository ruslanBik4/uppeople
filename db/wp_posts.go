// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_posts struct {
	dbEngine.Table
	Record *Wp_postsFields
	rows   sql.Rows
}

type Wp_postsFields struct {
	Id                    int64     `json:"id"`
	Post_author           float64   `json:"post_author"`
	Post_date             time.Time `json:"post_date"`
	Post_date_gmt         time.Time `json:"post_date_gmt"`
	Post_content          string    `json:"post_content"`
	Post_title            string    `json:"post_title"`
	Post_excerpt          string    `json:"post_excerpt"`
	Post_status           string    `json:"post_status"`
	Comment_status        string    `json:"comment_status"`
	Ping_status           string    `json:"ping_status"`
	Post_password         string    `json:"post_password"`
	Post_name             string    `json:"post_name"`
	To_ping               string    `json:"to_ping"`
	Pinged                string    `json:"pinged"`
	Post_modified         time.Time `json:"post_modified"`
	Post_modified_gmt     time.Time `json:"post_modified_gmt"`
	Post_content_filtered string    `json:"post_content_filtered"`
	Post_parent           float64   `json:"post_parent"`
	Guid                  string    `json:"guid"`
	Menu_order            int64     `json:"menu_order"`
	Post_type             string    `json:"post_type"`
	Post_mime_type        string    `json:"post_mime_type"`
	Comment_count         int64     `json:"comment_count"`
}

func (r *Wp_postsFields) RefColValue(name string) interface{} {
	switch name {
	case "id":
		return &r.Id

	case "post_author":
		return &r.Post_author

	case "post_date":
		return &r.Post_date

	case "post_date_gmt":
		return &r.Post_date_gmt

	case "post_content":
		return &r.Post_content

	case "post_title":
		return &r.Post_title

	case "post_excerpt":
		return &r.Post_excerpt

	case "post_status":
		return &r.Post_status

	case "comment_status":
		return &r.Comment_status

	case "ping_status":
		return &r.Ping_status

	case "post_password":
		return &r.Post_password

	case "post_name":
		return &r.Post_name

	case "to_ping":
		return &r.To_ping

	case "pinged":
		return &r.Pinged

	case "post_modified":
		return &r.Post_modified

	case "post_modified_gmt":
		return &r.Post_modified_gmt

	case "post_content_filtered":
		return &r.Post_content_filtered

	case "post_parent":
		return &r.Post_parent

	case "guid":
		return &r.Guid

	case "menu_order":
		return &r.Menu_order

	case "post_type":
		return &r.Post_type

	case "post_mime_type":
		return &r.Post_mime_type

	case "comment_count":
		return &r.Comment_count

	default:
		return nil
	}
}

func (r *Wp_postsFields) ColValue(name string) interface{} {
	switch name {
	case "id":
		return r.Id

	case "post_author":
		return r.Post_author

	case "post_date":
		return r.Post_date

	case "post_date_gmt":
		return r.Post_date_gmt

	case "post_content":
		return r.Post_content

	case "post_title":
		return r.Post_title

	case "post_excerpt":
		return r.Post_excerpt

	case "post_status":
		return r.Post_status

	case "comment_status":
		return r.Comment_status

	case "ping_status":
		return r.Ping_status

	case "post_password":
		return r.Post_password

	case "post_name":
		return r.Post_name

	case "to_ping":
		return r.To_ping

	case "pinged":
		return r.Pinged

	case "post_modified":
		return r.Post_modified

	case "post_modified_gmt":
		return r.Post_modified_gmt

	case "post_content_filtered":
		return r.Post_content_filtered

	case "post_parent":
		return r.Post_parent

	case "guid":
		return r.Guid

	case "menu_order":
		return r.Menu_order

	case "post_type":
		return r.Post_type

	case "post_mime_type":
		return r.Post_mime_type

	case "comment_count":
		return r.Comment_count

	default:
		return nil
	}
}

func NewWp_posts(db *dbEngine.DB) (*Wp_posts, error) {
	table, ok := db.Tables["wp_posts"]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: "wp_posts"}
	}

	return &Wp_posts{
		Table: table,
	}, nil
}

func (t *Wp_posts) NewRecord() *Wp_postsFields {
	t.Record = &Wp_postsFields{}
	return t.Record
}

func (t *Wp_posts) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_posts) SelectSelfScanEach(ctx context.Context, each func(record *Wp_postsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_posts) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_posts) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
