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
	"github.com/ruslanBik4/httpgo/auth"
	"github.com/ruslanBik4/logs"
	"golang.org/x/net/context"
)

var LastErr *pgconn.PgError
var initCustomTypes bool

type CitextArray struct {
	pgtype.TextArray
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

func GetDB(ctxApis apis.CtxApis) *dbEngine.DB {
	conn := psql.NewConn(AfterConnect, BeforeAcquire, printNotice)
	ctx := context.WithValue(ctxApis, "dbURL", "")
	ctx = context.WithValue(ctx, "fillSchema", true)
	db, err := dbEngine.NewDB(ctx, conn)
	if err != nil {
		logs.ErrorLog(err, "")
		return nil
	}

	err = initTagIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init TagIds")
	}

	err = initStatusesIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init StatusesIds")
	}

	err = initStatusesForVacIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init StatusesForVacsIds")
	}

	err = initSeniorityIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init SeniorityIds")
	}

	err = initPlatformIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init PlatformIds")
	}

	err = initLanguagesIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init LanguagesIds")
	}

	err = initLogActionsIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init LogActionsIds")
	}

	err = initLocationsIds(ctx, db)
	if err != nil {
		logs.ErrorLog(err, "on init LocationsIds")
	}

	LogsTable, err = NewLogs(db)
	if err != nil {
		logs.ErrorLog(err, "on init LogsTable")
	}

	return db
}

func BeforeAcquire(ctx context.Context, conn *pgx.Conn) bool {
	schema := GetSchema(ctx)
	if schema == "" {
		schema = "public"
	}

	tag, err := conn.Exec(ctx, "SET search_path TO "+schema)
	if err != nil {
		logs.ErrorLog(err, "SET search_path TO "+schema+tag.String())
		return false
	}

	return true
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
		mess += fmt.Sprintf("(%s, %v, %T) ", name, val.OID, val.Value)
	}

	logs.StatusLog(conn.PgConn().Conn().LocalAddr().String(), mess)

	return nil
}

func GetSchema(ctx context.Context) string {
	token, ok := ctx.Value(auth.UserValueToken).(interface{ GetSchema() string })
	if ok {
		return token.GetSchema()
	}

	return ""
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

func getOidCustomTypes(ctx context.Context, conn *pgx.Conn) error {
	params := make([]string, 0, len(customTypes))
	for name := range customTypes {
		params = append(params, name)
	}

	rows, err := conn.Query(ctx, SQL_GET_TYPES, params)
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
