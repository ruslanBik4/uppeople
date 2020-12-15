// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/httpgo/views"
	"github.com/ruslanBik4/httpgo/views/templates/layouts"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/data"
	"github.com/ruslanBik4/uppeople/views/pages"
)

// HandleLogServer show status httpgo
// @/api/version/
func HandleIndex(ctx *fasthttp.RequestCtx) (interface{}, error) {

	fWeb, ok := ctx.UserValue(WEB_PATH).(string)
	if !ok {
		return apis.ErrorResp{map[string]string{WEB_PATH: "is wrong"}}, apis.ErrRouteForbidden
	}

	filename := strings.TrimLeft(string(ctx.Request.URI().Path()), "/")
	b := bytes.Split(ctx.Request.URI().Host(), []byte("."))
	subDomain := string(b[0])
	if subDomain != "admin" {
		if filename == "" {
			filename = "index.html"
		}

		if len(b) < 3 {
			subDomain = "www/build"
		}

		fullName := filepath.Join(fWeb, subDomain, filename)
		_, err := os.Stat(fullName)
		if os.IsNotExist(err) {
			fullName = filepath.Join(fWeb, subDomain, "index.html")
		}

		if !bytes.Contains(ctx.Request.Host(), []byte("localhost")) {
			fasthttp.ServeFile(ctx, fullName)
			return nil, nil
		}

		b, err := ioutil.ReadFile(fullName)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return nil, nil
		}

		pageHtml := string(b)
		pageHtml = strings.ReplaceAll(pageHtml, "{{filepath}}", filepath.Join(fWeb, subDomain, filename))
		ctx.SetContentType("text/html; charset=utf-8")
		views.WriteHeadersHTML(ctx)
		ctx.Request.SetRequestURI(fullName)

		return pageHtml, nil
	}

	var err error
	if filename > "" {
		fullName := filepath.Join(fWeb, filename)
		_, err = os.Stat(fullName)
		if err == nil {

			fasthttp.ServeFile(ctx, fullName)

			return nil, nil
		}
	}

	body := &pages.IndexPageBody{
		TopMenu: []layouts.ItemMenu{
			{Label: "Translators", Link: path.Join(data.PathVersion, "table/dictionary/browse")},
			{Link: "#", Label: "New"},
			{Link: "#", Label: "View"},
			{Link: "#", Label: "Browse"},
			{Link: "#", Label: "Custom forms"},
		},
		HeadHTML: &layouts.HeadHTMLPage{
			Charset:  "charset=utf-8",
			Language: "ua",
			Title:    "Admin only! Save ours polymers!!!",
		},
		Title: "База даних сучасних полімерних матеріалів",
	}

	if err != nil {
		body.Content = fmt.Sprintf("Error on read file '%s': %s", filename, err.Error())
	}
	menuForms := layouts.Menu{}
	menuViews := layouts.Menu{}
	menuBrowse := layouts.Menu{}
	menuOther := layouts.Menu{}
	names := make([]string, 0, len(data.GlobalFormsList.Routes))
	for name := range data.GlobalFormsList.Routes {
		names = append(names, name)
	}

	sort.Strings(names)
	for _, name := range names {

		route := data.GlobalFormsList.Routes[name]
		words := strings.Split(name, "/")
		item := layouts.ItemMenu{
			Link:  name,
			Label: strings.Title(words[len(words)-2]),
			Title: route.Desc,
		}
		switch {
		case strings.HasSuffix(name, "form"):
			menuForms = append(menuForms, item)
		case strings.HasSuffix(name, "view"):
			item.Link += "?counter=100"
			menuViews = append(menuViews, item)
		case strings.HasSuffix(name, "browse"):
			item.Link += "?counter=100"
			menuBrowse = append(menuBrowse, item)
		case item.Label == "Forms":
			item.Label = strings.Title(words[len(words)-1])
			menuOther = append(menuOther, item)
		}
	}

	body.TopMenu[1].Content = menuForms.RenderDropdownMenu()
	body.TopMenu[2].Content = menuViews.RenderDropdownMenu()
	body.TopMenu[3].Content = menuBrowse.RenderDropdownMenu()
	body.TopMenu[4].Content = menuOther.RenderDropdownMenu()

	body.OwnerMenu = &layouts.MenuOwnerBody{
		TopMenu: layouts.Menu{
			layouts.ItemMenu{
				Link:    "/user/profile/",
				Label:   "Profile",
				Content: "",
				Title:   "",
			},
			// layouts.ItemMenu{
			// 	Link:    "#",
			// 	Label:   "Вхідні",
			// 	Content: "",
			// 	Title:   "",
			// },
			// layouts.ItemMenu{
			// 	Link:    "#",
			// 	Label:   "Запити",
			// 	Content: "",
			// 	Title:   "",
			// },
			// layouts.ItemMenu{
			// 	Link:    "#",
			// 	Label:   "Непрочитано",
			// 	Content: "",
			// 	Title:   "",
			// },
			// layouts.ItemMenu{
			// 	Link:    "#",
			// 	Label:   "Транспортуються",
			// 	Content: "",
			// 	Title:   "",
			// },
			// layouts.ItemMenu{
			// 	Link:    "#",
			// 	Label:   "В обробці",
			// 	Content: "",
			// 	Title:   "",
			// },
			// layouts.ItemMenu{
			// 	Link:    "#",
			// 	Label:   "Обрані",
			// 	Content: "",
			// 	Title:   "",
			// },
		},
		Title: "Kабінет:",
	}

	views.RenderHTMLPage(ctx, body.WriteShowMain)

	return nil, nil
}

var regTitle = regexp.MustCompile(`<title>([\W\w\s]+)</title>`)
