// generate file
// don't edit
package db

import (
	"database/sql"
	"time"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"golang.org/x/net/context"
)

type Wp_comments struct {
	dbEngine.Table
	Record *Wp_commentsFields
	rows   sql.Rows
}

type Wp_commentsFields struct {
	Comment_id           int64     `json:"comment_id"`
	Comment_post_id      float64   `json:"comment_post_id"`
	Comment_author       string    `json:"comment_author"`
	Comment_author_email string    `json:"comment_author_email"`
	Comment_author_url   string    `json:"comment_author_url"`
	Comment_author_ip    string    `json:"comment_author_ip"`
	Comment_date         time.Time `json:"comment_date"`
	Comment_date_gmt     time.Time `json:"comment_date_gmt"`
	Comment_content      string    `json:"comment_content"`
	Comment_karma        int64     `json:"comment_karma"`
	Comment_approved     string    `json:"comment_approved"`
	Comment_agent        string    `json:"comment_agent"`
	Comment_type         string    `json:"comment_type"`
	Comment_parent       float64   `json:"comment_parent"`
	User_id              float64   `json:"user_id"`
}

func (r *Wp_commentsFields) RefColValue(name string) interface{} {
	switch name {
	case "comment_id":
		return &r.Comment_id

	case "comment_post_id":
		return &r.Comment_post_id

	case "comment_author":
		return &r.Comment_author

	case "comment_author_email":
		return &r.Comment_author_email

	case "comment_author_url":
		return &r.Comment_author_url

	case "comment_author_ip":
		return &r.Comment_author_ip

	case "comment_date":
		return &r.Comment_date

	case "comment_date_gmt":
		return &r.Comment_date_gmt

	case "comment_content":
		return &r.Comment_content

	case "comment_karma":
		return &r.Comment_karma

	case "comment_approved":
		return &r.Comment_approved

	case "comment_agent":
		return &r.Comment_agent

	case "comment_type":
		return &r.Comment_type

	case "comment_parent":
		return &r.Comment_parent

	case "user_id":
		return &r.User_id

	default:
		return nil
	}
}

func (r *Wp_commentsFields) ColValue(name string) interface{} {
	switch name {
	case "comment_id":
		return r.Comment_id

	case "comment_post_id":
		return r.Comment_post_id

	case "comment_author":
		return r.Comment_author

	case "comment_author_email":
		return r.Comment_author_email

	case "comment_author_url":
		return r.Comment_author_url

	case "comment_author_ip":
		return r.Comment_author_ip

	case "comment_date":
		return r.Comment_date

	case "comment_date_gmt":
		return r.Comment_date_gmt

	case "comment_content":
		return r.Comment_content

	case "comment_karma":
		return r.Comment_karma

	case "comment_approved":
		return r.Comment_approved

	case "comment_agent":
		return r.Comment_agent

	case "comment_type":
		return r.Comment_type

	case "comment_parent":
		return r.Comment_parent

	case "user_id":
		return r.User_id

	default:
		return nil
	}
}

func NewWp_comments(db *dbEngine.DB) (*Wp_comments, error) {
	table, ok := db.Tables[TABLE_WPComments]
	if !ok {
		return nil, dbEngine.ErrNotFoundTable{Table: TABLE_WPComments}
	}

	return &Wp_comments{
		Table: table,
	}, nil
}

func (t *Wp_comments) NewRecord() *Wp_commentsFields {
	t.Record = &Wp_commentsFields{}
	return t.Record
}

func (t *Wp_comments) GetFields(columns []dbEngine.Column) []interface{} {
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

func (t *Wp_comments) SelectSelfScanEach(ctx context.Context, each func(record *Wp_commentsFields) error, Options ...dbEngine.BuildSqlOptions) error {
	return t.SelectAndScanEach(ctx,
		func() error {
			if each != nil {
				return each(t.Record)
			}

			return nil
		}, t, Options...)
}

func (t *Wp_comments) Insert(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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

func (t *Wp_comments) Update(ctx context.Context, Options ...dbEngine.BuildSqlOptions) (int64, error) {
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
