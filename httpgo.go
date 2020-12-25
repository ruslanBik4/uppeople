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
	"time"

	"github.com/ruslanBik4/httpgo/apis"
	httpgo "github.com/ruslanBik4/httpgo/httpGo"
	"github.com/ruslanBik4/httpgo/services"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"

	"github.com/ruslanBik4/uppeople/api"
	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
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
	DB := db.GetDB()
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
