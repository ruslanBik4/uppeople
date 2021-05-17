// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path"
	"regexp"
	"runtime/trace"
	"time"

	"github.com/ruslanBik4/httpgo/apis"
	httpgo "github.com/ruslanBik4/httpgo/httpGo"
	"github.com/ruslanBik4/httpgo/models/telegrambot"
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
	}
	fPort     = flag.String("port", ":443", "host address to listen on")
	fNoSecure = flag.Bool("insecure", false, "flag to force https")
	fSystem   = flag.String("path", "./", "path to system files")
	fCfgPath  = flag.String("config_path", "cfg", "path to cfg files")
	fWeb      = flag.String("web", "./", "path to web files")
)

// instances
var (
	httpServer *httpgo.HttpGo
)

func init() {
	flag.Parse()
	listener, err := net.Listen("tcp", *fPort)
	if err != nil {
		// port is occupied - work serve unpossable
		logs.Fatal(err)
	}

	ctxApis := apis.NewCtxApis(0)

	ctxApis.AddValue("migration", path.Join(*fCfgPath, "DB"))
	DB := db.GetDB(ctxApis)
	if DB == nil {
		panic("cannot init DB")
	}

	ctxApis.AddValue("DB", DB)
	ctxApis.AddValue(api.CFG_PATH, *fCfgPath)
	ctxApis.AddValue(api.SYSTEM_PATH, *fSystem)
	ctxApis.AddValue(api.WEB_PATH, *fWeb)
	ctxApis.AddValue("auth", auth.Bearer)
	ctxApis.AddValue("startedAt", time.Now())

	a := apis.NewApis(ctxApis, api.Routes, auth.Bearer)
	badRoutings := a.AddRoutes(routes)
	if len(badRoutings) > 0 {
		logs.ErrorLog(apis.ErrRouteForbidden, badRoutings)
	}

	cfg, err := httpgo.NewCfgHttp(path.Join(*fSystem, *fCfgPath, "httpgo.yml"))
	if err != nil || cfg == nil {
		// not work without correct config
		logs.Fatal(err, cfg)
	}

	httpServer = httpgo.NewHttpgo(cfg, listener, a)

	ctx := context.WithValue(context.TODO(), "mapRouting", api.Routes)
	services.InitServices(ctx, "mail", "showLogs")

}

var regIp = regexp.MustCompile(`for=s*(\d+\.?)+,`)

// version
var (
	Version string
	Build   string
	Branch  string
)

// HandleVersion show status httpgo
// @/api/version/
func HandleVersion(ctx *fasthttp.RequestCtx) (interface{}, error) {

	return map[string]interface{}{
		"UPPeople":  Branch,
		"Version":   Version,
		"BuildTime": Build,
		"StartTime": ctx.Value("startedAt"),
	}, nil
}

func main() {
	title := getAppTitle()

	t := "https"
	if *fNoSecure {
		t = "http"
	}

	logs.StatusLog("%s starting %s on port %s (%s)", title, time.Now(), *fPort, t)

	defer func() {
		errRec := recover()
		if err, ok := errRec.(error); ok {
			logs.ErrorStack(err)
		}
	}()

	ch := make(chan string)
	go func() {
		tBot, err := telegrambot.NewTelegramBotFromEnv()
		if err != nil {
			logs.ErrorLog(err, "NewTelegramBotFromEnv")
			return
		}
		logs.SetWriters(tBot, logs.FgErr, logs.FgDebug)
		if Branch > "" {
			logs.DebugLog(title)
			err, resp := tBot.SendMessage(title+"#starting", true)
			if err != nil {
				logs.ErrorLog(err, resp)
			}
		}

		msg := <-ch
		logs.DebugLog(msg)
		err, resp := tBot.SendMessage(
			fmt.Sprintf("#shutdown at %v %s", time.Now(), msg),
			true)
		if err != nil {
			logs.ErrorLog(err, resp)
		}
	}()

	if f, err := os.Create(Branch + ".out"); err != nil {
		logs.ErrorLog(err, "trace")
	} else {
		err = trace.Start(f)
		if err != nil {
			logs.ErrorLog(err, "trace")
			runServer()
		} else {
			defer trace.Stop()
			ctx, task := trace.NewTask(context.TODO(), "test")
			defer task.End()
			reg := trace.StartRegion(ctx, "httpgo")
			defer reg.End()
			logs.DebugLog(reg, task)
			trace.WithRegion(ctx, "httpgo", runServer)
		}
	}

	ch <- "finish"
	time.After(time.Second)
}

func getAppTitle() string {
	return fmt.Sprintf("UPPeople (%s) Version: %s, Build Time: %s", Branch, Version, Build)
}

func runServer() {
	err := httpServer.Run(
		!(*fNoSecure),
		path.Join(*fSystem, *fCfgPath, "server.crt"),
		path.Join(*fSystem, *fCfgPath, "server.key"))
	if err != nil {
		logs.ErrorStack(err)
	} else {
		logs.StatusLog("Server https correct shutdown at %v", time.Now())
	}

}
