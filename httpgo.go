// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/dbEngine/dbEngine/psql"
	"github.com/ruslanBik4/httpgo/apis"
	httpgo "github.com/ruslanBik4/httpgo/httpGo"
	"github.com/ruslanBik4/httpgo/services"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"

	"github.com/ruslanBik4/uppeople/api"
	"github.com/ruslanBik4/uppeople/auth"
)

//go:generate qtc -dir=views

const ShowVersion = "/api/version()"

var (
	routes = apis.ApiRoutes{
		"/": &apis.ApiRoute{
			Desc: "default endpoint",
			Fnc:  api.HandleIndex,
			// FncAuth: auth.Basic,
		},
		ShowVersion: {
			Fnc:  HandleVersion,
			Desc: "view version server",
		},
		// "/test/": &apis.ApiRoute{
		// 	Desc:   "default endpoint",
		// 	Fnc:    HandleTest,
		// 	Method: apis.POST,
		// },
	}
	fPort     = flag.String("port", ":443", "host address to listen on")
	fPortRdr  = flag.String("port_redirect", ":80", "redirect anather proto")
	fNoSecure = flag.Bool("insecure", false, "flag to force https")
	fSystem   = flag.String("path", "./", "path to system files")
	fCfgPath  = flag.String("config_path", "cfg", "path to cfg files")
	fWeb      = flag.String("web", "./", "path to web files")
)

var httpServer *httpgo.HttpGo

func init() {
	flag.Parse()
	listener, err := net.Listen("tcp", *fPort)
	if err != nil {
		// port is occupied - work serve unpossable
		logs.Fatal(err)
	}

	ctxApis := apis.NewCtxApis(0)

	ctxApis.AddValue(api.CFG_PATH, *fCfgPath)
	ctxApis.AddValue(api.SYSTEM_PATH, *fSystem)
	ctxApis.AddValue(api.WEB_PATH, *fWeb)
	DB := getDB()
	if DB == nil {
		logs.DebugLog(" ot DB")
	}

	ctxApis.AddValue("DB", DB)
	ctxApis.AddValue("auth", auth.Bearer)

	a := apis.NewApis(ctxApis, api.ApiRoutes, auth.Bearer)
	badRoutings := a.AddRoutes(routes)
	if len(badRoutings) > 0 {
		logs.ErrorLog(apis.ErrRouteForbidden, badRoutings)
	}

	// badRoutings = a.AddRoutes(auth2.AuthRoutes)
	// if len(badRoutings) > 0 {
	// 	logs.ErrorLog(apis.ErrRouteForbidden, badRoutings)
	// }
	//
	// badRoutings = a.AddRoutes(data.RoutesFromDB(DB, data.PathVersion))
	// if len(badRoutings) > 0 {
	// 	logs.ErrorLog(apis.ErrRouteForbidden, badRoutings)
	// }

	cfg, err := httpgo.NewCfgHttp(path.Join(*fSystem, *fCfgPath, "httpgo.yml"))
	if err != nil || cfg == nil {
		// not work without correct config
		logs.Fatal(err, cfg)
	}

	httpServer = httpgo.NewHttpgo(cfg, listener, a)

	ctx := context.WithValue(context.TODO(), "mapRouting", api.ApiRoutes)
	services.InitServices(ctx, "mail", "crypto", "showLogs")

	if ln, err := net.Listen("tcp", *fPortRdr); err != nil {
		// port is occupied - work without redirect
		logs.ErrorLog(err)
	} else {
		go func() {
			logs.StatusLog("Redirect service starting %s on port %s", time.Now(), *fPortRdr)
			err = fasthttp.Serve(ln, func(ctx *fasthttp.RequestCtx) {
				uri := ctx.Request.URI()
				proto := map[bool]string{
					true:  "http",
					false: "https",
				}
				uri.SetScheme(proto[*fNoSecure])
				if h := bytes.Split(uri.Host(), []byte(":")); len(h) > 1 {
					uri.SetHostBytes(h[0])
				}

				ctx.RedirectBytes(uri.FullURI(), fasthttp.StatusMovedPermanently)
				logs.DebugLog("redirect %s", string(uri.FullURI()))
				logIP(ctx)

			})
			if err != nil {
				logs.ErrorLog(err, "fasthttpServe")
			}

		}()
	}
}

var regIp = regexp.MustCompile(`for=s*(\d+\.?)+,`)

func logIP(ctx *fasthttp.RequestCtx) {
	ipClient := ctx.Request.Header.Peek("X-Forwarded-For")
	addr := string(ipClient)
	if len(ipClient) == 0 {
		ipClient = ctx.Request.Header.Peek("Forwarded")
		ips := regIp.FindSubmatch(ipClient)

		if len(ips) == 0 {
			addr = string(ctx.Request.Header.Peek("X-ProxyUser-Ip"))
		} else {
			addr = string(ips[0])
		}
	}

	if addr == "" {
		logs.StatusLog(ctx.RemoteAddr().String())
	} else {
		logs.StatusLog(addr)
	}
}

func getDB() *dbEngine.DB {
	conn := psql.NewConn(AfterConnect, nil, printNotice)
	ctx := context.WithValue(context.Background(), "dbURL", "")
	ctx = context.WithValue(ctx, "fillSchema", true)
	ctx = context.WithValue(ctx, "migration", "cfg/DB")
	db, err := dbEngine.NewDB(ctx, conn)
	if err != nil {
		logs.ErrorLog(err, "")
		return nil
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
		logs.ErrorLog(err, n.Hint, err.SQLState(), err.File, err.Line, err.Routine)
	} else if strings.HasPrefix(n.Message, "[[ERROR]]") {
		logs.ErrorLog(errors.New(strings.TrimPrefix(n.Message, "[[ERROR]]") + n.Severity))
	} else { // DEBUG
		logs.DebugLog("%+v %s", n.Severity, n.Message)
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

// func (dst *CitextArray) MarshalJSON() ([]byte, error) {
// 	buf := bytes.NewBufferString("[")
// 	for i, text := range dst.Elements {
// 		if i > 0 {
// 			buf.WriteString(",")
// 		}
//
// 		buf.WriteString(text.String)
// 	}
//
// 	buf.WriteString("]")
//
// 	return buf.Bytes(), nil
// }

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
	"citext": &pgtype.DataType{
		Value: &pgtype.Text{},
		Name:  "citext",
	},
	"_citext": &pgtype.DataType{
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

// version
var (
	Version string
	Build   string
	Branch  string
)

// HandleLogServer show status httpgo
// @/api/version/
func HandleVersion(ctx *fasthttp.RequestCtx) (interface{}, error) {

	return fmt.Sprintf("UPPeople (%s) Version: %s, Build Time: %s", Branch, Version, Build), nil
}

func main() {
	title, err := HandleVersion(nil)

	t := "https"
	if *fNoSecure {
		t = "http"
	}

	logs.StatusLog("%s starting %s on port %s (%s)", title, time.Now(), *fPort, t)

	defer func() {
		errRec := recover()
		if err, ok := errRec.(error); ok {
			logs.ErrorLog(err)
		}
	}()

	err = httpServer.Run(
		!(*fNoSecure),
		path.Join(*fSystem, *fCfgPath, "server.crt"),
		path.Join(*fSystem, *fCfgPath, "server.key"))
	if err != nil {
		logs.ErrorLog(err)
	} else {
		logs.StatusLog("Server https correct shutdown at %v", time.Now())
	}
}
