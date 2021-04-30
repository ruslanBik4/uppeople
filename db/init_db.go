// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package db

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/dbEngine/dbEngine/psql"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"golang.org/x/net/context"
)

var LastErr *pgconn.PgError

func GetDB(ctxApis apis.CtxApis) *dbEngine.DB {
	conn := psql.NewConn(AfterConnect, nil, printNotice)
	ctx := context.WithValue(ctxApis, "dbURL", "")
	ctx = context.WithValue(ctx, "fillSchema", true)
	db, err := dbEngine.NewDB(ctx, conn)
	if err != nil {
		logs.ErrorLog(err, "")
		return nil
	}

	tagIds = &TagIdMap{}
	tagsTable, err := NewTags(db)
	if err != nil {
		logs.ErrorLog(err, "cannot get %s table", TableTags)
		return db
	}

	err = tagsTable.SelectSelfScanEach(ctx,
		func(record *TagsFields) error {
			(*tagIds)[record.Name] = *record
			return nil
		})

	if err != nil {
		logs.ErrorLog(err, "while reading tags from db to tagIds(db.TagIdMap)")
	}

	return db
}

func printNotice(c *pgconn.PgConn, n *pgconn.Notice) {

	if n.Code == "42P07" || strings.Contains(n.Message, "skipping") {
		logs.DebugLog("skip operation: %s", n.Message)
	} else if n.Severity == "INFO" {
		logs.StatusLog(n.Message)
	} else if n.Code > "00000" {
		err := (*pgconn.PgError)(n)
		LastErr = err
		logs.ErrorLog(err, n.Hint, err.SQLState(), err.File, err.Line, err.Routine)
	} else if strings.HasPrefix(n.Message, "[[ERROR]]") {
		logs.ErrorLog(errors.New(strings.TrimPrefix(n.Message, "[[ERROR]]") + n.Severity))
	} else { // DEBUG
		logs.DebugLog("%d: %+v %s", c.PID(), n.Severity, n.Message)
	}
}

var initCustomTypes bool

const sqlGetTypes = "SELECT typname, oid FROM pg_type WHERE typname::text=ANY($1)"

type CitextArray struct {
	pgtype.TextArray
}

func (dst CitextArray) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, text := range dst.Elements {
		if i > 0 {
			buf.WriteString(",")
		}

		buf.WriteString(text.String)
	}

	buf.WriteString("]")

	return buf.Bytes(), nil
}

func (dst CitextArray) Get1() interface{} {
	res := make([]string, len(dst.Elements))
	logs.DebugLog(dst)
	return &res
}
func (dst *CitextArray) Set(src interface{}) error {
	switch value := src.(type) {
	case []string:
		logs.DebugLog(src)
		elements := make([]pgtype.Text, len(value))
		for i, str := range value {
			elements[i].String = str
			elements[i].Status = pgtype.Present
			(*dst).TextArray = pgtype.TextArray{
				Elements:   elements,
				Dimensions: []pgtype.ArrayDimension{{Length: int32(len(elements)), LowerBound: 1}},
				Status:     pgtype.Present,
			}
		}
		return nil
	default:
		return dst.TextArray.Set(src)
	}
}
func (src CitextArray) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case pgtype.Null:
		return nil, nil
		// case pgtype.Undefined:
		// 	return nil, pgtype.ErrUndefined
	case pgtype.Present:

		buf = append(buf, '"', '{')
		for i, elem := range src.Elements {
			if i > 0 {
				buf = append(buf, ',')
			}

			buf = append(buf, pgtype.QuoteArrayElementIfNeeded(elem.String)...)
		}
		buf = append(buf, '}', '"')

		logs.DebugLog(string(buf))
		return buf, nil
	default:

		buf, err := src.TextArray.EncodeText(ci, buf)
		if err != nil {
			logs.ErrorLog(err)
			return nil, errors.Wrap(err, "")
		}

		logs.DebugLog(string(buf))
		return buf, nil
	}
}

func (src CitextArray) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	switch src.Status {
	case pgtype.Present:

		arrayHeader := pgtype.ArrayHeader{
			Dimensions: src.Dimensions,
		}

		if dt, ok := ci.DataTypeForName("citext"); ok {
			arrayHeader.ElementOID = int32(dt.OID)
		} else {
			return nil, errors.Errorf("unable to find oid for type name %v", "text")
		}

		for i := range src.Elements {
			if src.Elements[i].Status == pgtype.Null {
				arrayHeader.ContainsNull = true
				break
			}
		}

		buf = arrayHeader.EncodeBinary(ci, buf)

		for i := range src.Elements {
			sp := len(buf)
			buf = pgio.AppendInt32(buf, -1)

			elemBuf, err := src.Elements[i].EncodeBinary(ci, buf)
			if err != nil {
				return nil, err
			}
			if elemBuf != nil {
				buf = elemBuf
				pgio.SetInt32(buf[sp:], int32(len(buf[sp:])-4))
			}
		}

		return buf, nil
	default:
		return src.TextArray.EncodeBinary(ci, buf)
	}
}

var customTypes = map[string]*pgtype.DataType{
	"citext": {
		Value: &pgtype.Text{},
		Name:  "citext",
	},
	"_citext": {
		Value: &CitextArray{pgtype.TextArray{}}, // (*pgtype.ArrayType)(nil),
		Name:  "[]string",
	},
}

func AfterConnect(ctx context.Context, conn *pgx.Conn) error {
	// Override registered handler for point
	if !initCustomTypes {
		err := getOidCustomTypes(ctx, conn)
		if err != nil {
			return err
		}

		initCustomTypes = true
	}

	mess := "DB registered type (name, oid): "
	for name, val := range customTypes {
		conn.ConnInfo().RegisterDataType(*val)
		mess += fmt.Sprintf("(%s,%v, %T) ", name, val.OID, val.Value)
	}

	logs.StatusLog(conn.PgConn().Conn().LocalAddr().String(), mess)

	return nil
}
func getOidCustomTypes(ctx context.Context, conn *pgx.Conn) error {
	params := make([]string, 0, len(customTypes))
	for name := range customTypes {
		params = append(params, name)
	}

	rows, err := conn.Query(ctx, sqlGetTypes, params)
	if err != nil {
		return err
	}

	for rows.Next() {
		var name string
		var oid uint32
		err = rows.Scan(&name, &oid)
		if err != nil {
			return err
		}
		if c, ok := customTypes[name]; ok && c.Value == (*pgtype.ArrayType)(nil) {
			c.Value = pgtype.NewArrayType(name, oid, func() pgtype.ValueTranscoder {
				return &pgtype.Text{}
			}).NewTypeValue()
			c.OID = oid
			logs.DebugLog(c)
		} else if ok {
			customTypes[name].OID = oid
		}
	}

	if rows.Err() != nil {
		logs.ErrorLog(rows.Err(), " cannot get oid for customTypes")
	}

	return err
}
