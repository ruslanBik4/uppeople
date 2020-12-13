// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"bytes"
	"fmt"
	"go/types"
	"io/ioutil"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/httpgo/views"
	"github.com/ruslanBik4/httpgo/views/templates/forms"
	"github.com/ruslanBik4/httpgo/views/templates/json"
	"github.com/ruslanBik4/httpgo/views/templates/system/routeTable"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/auth"
	"github.com/ruslanBik4/uppeople/db"
	"github.com/ruslanBik4/uppeople/views/td"
)

func GetTable(ctx *fasthttp.RequestCtx, tableName string) (dbEngine.Table, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}
	table, ok := DB.Tables[tableName]
	if !ok {
		return nil, errors.New("table '" + tableName + "' not found")
	}

	return table, nil
}

var GlobalFormsList *FormsList

func RoutesFromDB(DB *dbEngine.DB, pathVersion string) apis.ApiRoutes {

	preRoute := pathVersion + "/table/"
	routes := make(apis.ApiRoutes, 0)
	GlobalFormsList = NewFormsList()
	patternList, err := db.NewPatterns_list(DB)
	if err != nil {
		logs.ErrorLog(err, "db.NewPatterns_list")
		return nil
	}

	for name, table := range DB.Tables {

		// if strings.HasPrefix(table.Name(), "form_") {
		// 	continue
		// }

		rUpd := &apis.ApiRoute{
			Desc:        "update table '" + name + "' data",
			Method:      apis.POST,
			Multipart:   true,
			DTO:         dtoField{},
			FncAuth:     nil,
			TestFncAuth: nil,
			NeedAuth:    true,
			OnlyAdmin:   false,
			OnlyLocal:   false,
			Params:      []apis.InParam{ParamsLang, ParamsGetFormActions},
			Resp:        nil,
		}

		rIns := &apis.ApiRoute{
			Desc:      "insert into table '" + name + "' data",
			Method:    apis.POST,
			Multipart: true,
			NeedAuth:  true,
			DTO:       dtoField{},
			Params:    []apis.InParam{ParamsLang, ParamsGetFormActions},
		}

		params := make([]string, 0)
		priColumns := make([]string, 0)
		basicParams := []apis.InParam{
			ParamsHTML,
			ParamsLang,
		}

		for _, col := range table.Columns() {

			p := apis.InParam{
				Name:              col.Name(),
				Desc:              col.Comment(),
				Req:               col.Primary(),
				PartReq:           nil,
				IncompatibleWiths: nil,
				TestValue:         "",
			}

			if strings.HasPrefix(col.Type(), "_") {
				p.Type = apis.NewSliceTypeInParam(col.BasicType())
				p.Name += "[]"
			} else if col.Type() == "date" {
				// todo add new type of date/time
				p.Type = apis.NewTypeInParam(types.String)
			} else {
				p.Type = apis.NewTypeInParam(col.BasicType())
			}

			rUpd.Params = append(rUpd.Params, p)

			if !col.AutoIncrement() {
				p.Req = col.Required()
				p.DefValue = col.Default()
				rIns.Params = append(rIns.Params, p)
				params = append(params, p.Name)
			}

			if col.Primary() || (col.Name() == "id") {
				priColumns = append(priColumns, p.Name)
				pForm := p
				pForm.Req = false
				basicParams = append(basicParams, pForm)
			}
		}

		// if !strings.HasPrefix(table.Name(), "form") {

		pathForm := preRoute + name + "/form"

		routes[pathForm] = &apis.ApiRoute{
			Desc:     "get form for insert/update data into " + name,
			Fnc:      TableForm(DB, preRoute, table, patternList, priColumns),
			NeedAuth: true,
			Params:   basicParams,
		}

		GlobalFormsList.Add(pathForm, name)

		rUpd.Fnc = TableUpdate(preRoute, table, params, priColumns)

		routes[preRoute+name+"/update"] = rUpd

		report := NewReportJSON(table)
		routes[preRoute+name+"/report"] = report.getRoute()
		routes[preRoute+name+"/data"] = &apis.ApiRoute{
			Desc: "from  table '" + table.Name() + "' data",
			Params: []apis.InParam{
				ParamsLang,
			},
			Fnc: TableData(DB, table, pathForm),
		}

		rIns.Fnc = TableInsert(preRoute, table, params)

		routes[preRoute+name+"/put"] = rIns

		// }

		routes[preRoute+name+"/view"] = &apis.ApiRoute{
			Desc: "view data of table " + table.Name(),
			Fnc:  TableView(preRoute, DB, table, patternList, priColumns),
			Params: []apis.InParam{
				ParamsLang,
				ParamsCounter,
			},
		}

		routes[preRoute+name+"/browse"] = &apis.ApiRoute{
			Desc:      "show & edit data of table " + table.Name(),
			NeedAuth:  true,
			OnlyAdmin: true,
			Fnc:       TableView(preRoute, DB, table, patternList, priColumns),
			Params: []apis.InParam{
				ParamsLang,
				ParamsCounter,
			},
		}

		routes[preRoute+name+"/"] = &apis.ApiRoute{
			Desc: "show row of table according to ID" + table.Name(),
			// NeedAuth: true,
			Fnc: TableRow(table),
			Params: []apis.InParam{
				ParamsLang,
				{
					Name: "id",
					Desc: "id of photos record for download",
					Req:  false,
					Type: apis.NewTypeInParam(types.Int32),
				},
			},
		}
	}

	err = GlobalFormsList.AddCustomForms(DB, routes, preRoute, patternList)
	if err != nil {
		logs.ErrorLog(err, "GlobalFormsList.AddCustomForms")
	}

	routes[pathVersion+"/forms/list/"] = GlobalFormsList.GetRoute()

	GlobalFormsList.Routes = routes

	return routes
}

func TableData(DB *dbEngine.DB, table dbEngine.Table, pathForm string) func(ctx *fasthttp.RequestCtx) (interface{}, error) {
	colSelect := make([]string, 0)
	for _, col := range table.Columns() {
		if col.BasicType() != types.UnsafePointer {
			colSelect = append(colSelect, col.Name())
		} else {
			colSelect = append(colSelect, `'null' :: bytea as `+col.Name())
		}
	}
	conn := DB.Conn
	sql := fmt.Sprintf("select %s from %s", strings.Join(colSelect, ","), table.Name())

	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {

		firstRows := ""
		if auth.Bearer.Auth(ctx) {
			user := auth.GetUserData(ctx)
			if user != nil {
				for id := range user.Companies {
					if firstRows == "" {
						firstRows = " where id=ANY(ARRAY["
					} else {
						firstRows += ","
					}
					firstRows += strconv.Itoa(int(id))
				}
			}

			if firstRows > "" {
				firstRows += "]) "
			}
			logs.DebugLog(user, firstRows)
		}

		data := make([]map[string]interface{}, 0)
		err := conn.SelectAndRunEach(ctx,
			func(row []interface{}, columns []dbEngine.Column) error {
				rec := make(map[string]interface{})
				id := int32(-1)
				for i, col := range columns {
					if col.Name() == "id" {
						id = row[i].(int32)
					}
					if v := GetForeignName(ctx, DB, col, row[i]); v != nil {
						rec[col.Name()] = v
					} else {
						rec[col.Name()] = ToStandartColumnValueType(table.Name(), col.Name(), id, row[i])
					}
				}

				// todo : check admin role
				// if auth.AuthBearer.AdminAuth(ctx) {
				rec["edit"] = fmt.Sprintf(pathForm+"?id=%d", id)
				// }

				data = append(data, rec)
				return nil
			},
			sql+firstRows)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
}

func TableForm(DB *dbEngine.DB, preRoute string, table, patternList dbEngine.Table, priColumns []string) apis.ApiRouteHandler {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		f := forms.FormField{
			Title:       table.Comment(),
			Action:      preRoute + table.Name(),
			Method:      "POST",
			Description: "Input data for " + table.Comment(),
		}

		// we must copy colsTable into local array
		colDecors := make([]*forms.ColumnDecor, 0)

		id, ok := int32(0), false
		args := make([]interface{}, len(priColumns))
		clsTable := table.Columns()
		for i, name := range priColumns {
			if name == "id" {
				id, ok = ctx.UserValue(name).(int32)
				if !ok {
					var s string
					s, ok = ctx.UserValue(apis.ChildRoutePath).(string)
					if ok {
						i, err := strconv.Atoi(s)
						if err != nil {
							return apis.ChildRoutePath, apis.ErrWrongParamsList
						}
						id = int32(i)
					}
				}
				args[i] = id
			} else {
				args[i] = ctx.UserValue(name)
				ok = args[i] != nil
			}

			if !ok {
				break
			}

			for _, col := range clsTable {
				if col.Name() == name {
					colDec := forms.NewColumnDecor(col, patternList)
					colDec.Value = args[i]
					if colDec.AutoIncrement() || col.Name() == "id" {
						colDec.IsHidden = true
						colDec.InputType = "hidden"
					} else if strings.Contains(col.Comment(), " (read_only)") {
						colDec.IsReadOnly = true
						colDec.IsDisabled = true
						colDec.Value = GetForeignName(ctx, DB, colDec, args[i])
					} else {
						// add hidden field for update where clause
						colOld := forms.NewColumnDecor(col, patternList)
						colOld.IsHidden = true
						colOld.InputType = "hidden"
						colOld.Value = args[i]
						colDecors = append(colDecors, colOld)

						GetForeignOptions(ctx, DB, colDec, args[i])
						colDec.IsNewPrimary = true
					}
					colDecors = append(colDecors, colDec)
					break
				}
			}
		}

		btnList := []forms.Button{
			{ButtonType: "submit", Title: "Insert", Position: true},
			{ButtonType: "reset", Title: "Clear", Position: false},
		}

		if ok {
			f.Action += "/update"
			btnList[0].Title = "Update"
			colSelect := make([]string, 0, len(clsTable)-len(priColumns))
		loop_columns:
			for _, colDec := range clsTable {
				for _, name := range priColumns {
					if name == colDec.Name() {
						continue loop_columns
					}
				}
				colSelect = append(colSelect, colDec.Name())
			}

			err := table.SelectAndRunEach(ctx,
				func(values []interface{}, columns []dbEngine.Column) error {
					ok = false
					for i, col := range columns {
						name := col.Name()
						values[i] = ToStandartColumnValueType(table.Name(), name, id, values[i])

						for _, col := range clsTable {
							if col.Name() == name {
								colDec := forms.NewColumnDecor(col, patternList)
								colDec.IsDisabled = colDec.IsReadOnly
								colDec.IsSlice = (strings.HasPrefix(col.Type(), "_"))

								if col.Type() == "text" {
									colDec.InputType = "textarea"
								}

								colDec.Value = values[i]
								GetForeignOptions(ctx, DB, colDec, values[i])
								colDecors = append(colDecors, colDec)
								break
							}
						}

					}

					return nil
				},
				dbEngine.ColumnsForSelect(colSelect...),
				dbEngine.WhereForSelect(priColumns...),
				dbEngine.ArgsForSelect(args...),
			)
			if err != nil {
				logs.ErrorLog(err, "")
			}

			// not found record
			if ok {
				ctx.SetStatusCode(fasthttp.StatusNoContent)
				return nil, nil
			}
		} else {
			f.Action += "/put"
			for _, col := range clsTable {
				if !(col.AutoIncrement() || col.Name() == "id" ||
					strings.Contains(col.Comment(), " (read_only)")) {

					colDec := forms.NewColumnDecor(col, patternList)
					colDec.Value = col.Default()
					if colDec.Value == "NULL" {
						colDec.Value = nil
					}

					colDec.IsSlice = (strings.HasPrefix(col.Type(), "_"))

					if col.Type() == "text" {
						colDec.InputType = "textarea"
					}
					GetForeignOptions(ctx, DB, colDec, nil)
					colDecors = append(colDecors, colDec)
				}

			}
		}

		lang, ok := ctx.UserValue(ParamsLang.Name).(string)
		if ok {
			colDecors = append(colDecors, &forms.ColumnDecor{
				Column:      dbEngine.NewStringColumn("lang", "lang", true),
				IsHidden:    true,
				InputType:   "hidden",
				PatternList: nil,
				Value:       lang,
			})
		}

		blocks := []forms.BlockColumns{
			{
				Buttons:     btnList,
				Columns:     colDecors,
				Id:          1,
				Title:       "",
				Description: "",
			},
		}

		_, ok = ctx.UserValue("html").(bool)
		if !ok {
			views.WriteJSONHeaders(ctx)
		}

		f.WriteRenderForm(
			ctx.Response.BodyWriter(),
			ok, // && isHtml,
			blocks...)

		return nil, nil
	}
}

func copyColumnDecor(ctx *fasthttp.RequestCtx, colDec *forms.ColumnDecor) *forms.ColumnDecor {
	c := colDec.Copy()
	if colDec.PlaceHolder > "" {
		c.PlaceHolder = Translate(ctx, colDec.PlaceHolder)
	}

	c.Label = Translate(ctx, colDec.Label)

	return c
}

func ToStandartColumnValueType(tableName, colName string, id int32, values interface{}) interface{} {
	// todo- move to dbEngine
	switch v := values.(type) {
	case pgtype.VarcharArray:
		return VarcharArrayToStrings(v.Elements)

	case *pgtype.VarcharArray:
		return VarcharArrayToStrings(v.Elements)

	case pgtype.TextArray:
		return TextArrayToStrings(v.Elements)

	case *pgtype.TextArray:
		return TextArrayToStrings(v.Elements)

	case pgtype.BPCharArray:
		return BPCharArrayToStrings(v.Elements)

	case *pgtype.BPCharArray:
		return BPCharArrayToStrings(v.Elements)

	case pgtype.Int4Array:
		return Int4ArrToStrings(v.Elements)

	case pgtype.Int8Array:
		return Int8ArrToStrings(v.Elements)

	case pgtype.ArrayType:
		str, done := ArrayToStrings(&v)
		if done {
			return str
		}

		return v

	case *pgtype.ArrayType:
		str, done := ArrayToStrings(v)
		if done {
			return str
		}

		return v

	case *pgtype.GenericText:
		logs.DebugLog("%T", v)
		return "genericText: " + v.String

	case pgtype.UntypedTextArray:
		return v.Elements

	case *pgtype.UntypedTextArray:
		return v.Elements

	case []interface{}:
		return UnknownArrayToStrings(v)

	case *pgtype.Bytea, pgtype.Bytea, []uint8:
		return BlobToURL(tableName, colName, id)

	case time.Time:
		return v.Format("2006-01-02")

	case *time.Time:
		return v.Format("2006-01-02")

	case nil, string, bool, float32, float64, int32, int64, map[string]string, map[string]interface{}:
		return values

	default:
		logs.DebugLog("%T", values)
		return values
	}
}

func BlobToURL(tableName string, colName string, id int32) string {
	return fmt.Sprintf("/api/v1/blob/%s?id=%d&name=%s", tableName, id, colName)
}

func ArrayToStrings(v *pgtype.ArrayType) ([]string, bool) {
	src, ok := v.Get().([]interface{})
	if !ok {
		return nil, false
	}

	return UnknownArrayToStrings(src), true
}

func Int4ArrToStrings(src []pgtype.Int4) []int32 {
	str := make([]int32, len(src))
	for i, val := range src {
		str[i] = val.Int
	}

	return str
}

func Int8ArrToStrings(src []pgtype.Int8) []int64 {
	str := make([]int64, len(src))
	for i, val := range src {
		str[i] = val.Int
	}

	return str
}

func UnknownArrayToStrings(src []interface{}) []string {
	str := make([]string, len(src))
	for i, val := range src {
		str[i] = json.Element(val)
	}

	return str
}

func VarcharArrayToStrings(src []pgtype.Varchar) []string {
	str := make([]string, len(src))
	for i, val := range src {
		str[i] = val.String
	}

	return str
}

func TextArrayToStrings(src []pgtype.Text) []string {
	str := make([]string, len(src))
	for i, val := range src {
		str[i] = val.String
	}

	return str
}

func BPCharArrayToStrings(src []pgtype.BPChar) []string {
	str := make([]string, len(src))
	for i, val := range src {
		str[i] = val.String
	}

	return str
}

func TableUpdate(preRoute string, table dbEngine.Table, columns, priColumns []string) apis.ApiRouteHandler {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {

		badParams := make(map[string]string, 0)
		for _, key := range priColumns {
			if ctx.UserValue(key) == nil {
				badParams[key] = "required params"
			}
		}

		if len(badParams) > 0 {
			return badParams, apis.ErrWrongParamsList

		}

		args := make([]interface{}, 0, len(columns))
		colSel := make([]string, 0, len(columns))
		msg := ""
		for _, name := range columns {
			isPrimary := false
			for _, priName := range priColumns {
				if name == priName {
					arg := ctx.UserValue("new." + priName)
					if arg != nil {
						colSel = append(colSel, priName)
						args = append(args, arg)
						msg += fmt.Sprintf(" %v", arg)
					}
					isPrimary = true
					break
				}
			}

			if isPrimary {
				continue
			}

			arg := ctx.UserValue(name)
			if arg == nil {
				continue
			}

			colName := strings.TrimSuffix(name, "[]")
			col := table.FindColumn(colName)
			if col.BasicType() == types.UnsafePointer {
				switch val := arg.(type) {
				case nil, string:
				case []*multipart.FileHeader:
					names, bytea, err := readByteA(val)
					if err != nil {
						logs.DebugLog(names)
						return map[string]string{colName: err.Error()}, apis.ErrWrongParamsList
					}

					switch len(bytea) {
					case 0:
					case 1:
						args = append(args, bytea[0])
						colSel = append(colSel, colName)
						msg += "[file]"
					default:
						args = append(args, bytea)
						colSel = append(colSel, colName)
						msg += "[files]"
					}
				default:
					return map[string]string{colName: fmt.Sprintf("%v", val)}, apis.ErrWrongParamsList
				}

				continue
			}

			args = append(args, arg)
			colSel = append(colSel, colName)
			msg += fmt.Sprintf(" %v", arg)

		}

		for _, name := range priColumns {
			args = append(args, ctx.UserValue(name))
		}

		i, err := table.Update(ctx,
			dbEngine.ColumnsForSelect(colSel...),
			dbEngine.WhereForSelect(priColumns...),
			dbEngine.ArgsForSelect(args...),
		)
		if err != nil {
			return nil, err
		}
		if i <= 0 {
			logs.DebugLog(colSel, priColumns)
			logs.DebugLog(args)
			return map[string]string{"update": fmt.Sprintf("%d", i)}, apis.ErrWrongParamsList
		}

		msg = "Success update: " + strings.Join(colSel, ", ") + " values:\n" + msg

		ctx.SetStatusCode(fasthttp.StatusAccepted)
		g, ok := ctx.UserValue(ParamsGetFormActions.Name).(bool)
		if ok && g {
			urlSuffix := "/browse"
			lang := ctx.UserValue("lang")
			if l, ok := lang.(string); ok {
				urlSuffix += "?lang=" + l
			}

			return insertResult{
				FormActions: []FormActions{
					{
						Typ: "redirect",
						Url: preRoute + table.Name() + urlSuffix,
					},
				},
				Msg: msg,
			}, nil
		}

		return msg, nil
	}
}

func readByteA(fHeaders []*multipart.FileHeader) ([]string, [][]byte, error) {
	bytea := make([][]byte, len(fHeaders))
	names := make([]string, len(fHeaders))
	for i, fHeader := range fHeaders {

		f, err := fHeader.Open()
		if err != nil {
			logs.DebugLog(err, fHeader)
			return nil, nil, errors.Wrap(err, fHeader.Filename)
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			logs.DebugLog(err, fHeader)
			return nil, nil, errors.Wrap(err, "read "+fHeader.Filename)
		}

		bytea[i] = b
		names[i] = fHeader.Filename
	}

	return names, bytea, nil
}

type insertResult struct {
	FormActions []FormActions `json:"formActions"`
	Msg         string        `json:"message"`
}

var regDublicate = regexp.MustCompile(`duplicate key value violates unique constraint "(\w*)"`)

// TableInsert insert data of params into table
func TableInsert(preRoute string, table dbEngine.Table, params []string) apis.ApiRouteHandler {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {

		args := make([]interface{}, 0, len(params))
		colSel := make([]string, 0, len(params))
		msg := ""
		for _, name := range params {
			arg := ctx.UserValue(name)
			if arg == nil {
				continue
			}

			colName := strings.TrimSuffix(name, "[]")
			col := table.FindColumn(colName)
			if col == nil {
				logs.ErrorLog(dbEngine.ErrNotFoundColumn{Table: table.Name(), Column: colName})
				continue
			}

			if col.BasicType() == types.UnsafePointer {

				switch val := arg.(type) {
				case nil, string:
				case []*multipart.FileHeader:

					names, bytea, err := readByteA(val)
					if err != nil {
						logs.DebugLog(names)
						return map[string]string{colName: err.Error()},
							apis.ErrWrongParamsList
					}

					switch len(bytea) {
					case 0:
					case 1:
						args = append(args, bytea[0])
						colSel = append(colSel, colName)
						msg += "[file]"
					default:
						args = append(args, bytea)
						colSel = append(colSel, colName)
						msg += "[files]"
					}
				default:
					return map[string]string{colName: fmt.Sprintf("%v", val)},
						apis.ErrWrongParamsList
				}

				continue
			}

			args = append(args, arg)
			colSel = append(colSel, colName)
			msg += fmt.Sprintf(" %v", arg)

		}

		_, err := table.Insert(ctx,
			dbEngine.ColumnsForSelect(colSel...),
			dbEngine.ArgsForSelect(args...),
		)
		if err != nil {
			e, ok := errors.Cause(err).(*pgconn.PgError)
			if ok {
				reg := regexp.MustCompile(`Key\s+\((\w+)\)=\((\w+)\)([^.]+)`)
				if s := reg.FindStringSubmatch(e.Detail); len(s) > 0 {
					return map[string]string{
						s[1]: "`" + s[2] + "`" + s[3],
					}, apis.ErrWrongParamsList
				}
			} else if s := regDublicate.FindStringSubmatch(err.Error()); len(s) > 0 {
				logs.DebugLog("%#v %[1]T", errors.Cause(err))
				return map[string]string{
					s[1]: "duplicate key value violates unique constraint",
				}, apis.ErrWrongParamsList
			}

			return nil, err
		}

		msg = "Success saving: " + strings.Join(colSel, ", ") + " values:\n" + msg

		ctx.SetStatusCode(fasthttp.StatusCreated)
		g, ok := ctx.UserValue(ParamsGetFormActions.Name).(bool)
		if ok && g {
			urlSuffix := "/form?html"

			lang := ctx.UserValue("lang")
			if l, ok := lang.(string); ok {
				urlSuffix += "&lang=" + l
			}

			return insertResult{
				FormActions: []FormActions{
					{
						Typ: "redirect",
						Url: preRoute + table.Name() + urlSuffix,
					},
				},
				Msg: msg,
			}, nil
		}

		return msg, nil
	}
}

// TableRow return one record of table according to id key
func TableRow(table dbEngine.Table) apis.ApiRouteHandler {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		rows := make(map[string]interface{}, len(table.Columns()))

		id, ok := ctx.UserValue("id").(int32)
		if !ok {
			s, ok := ctx.UserValue(apis.ChildRoutePath).(string)
			if !ok {
				return map[string]string{"id": "wrong"}, apis.ErrWrongParamsList
			}

			i, err := strconv.Atoi(s)
			if err != nil {
				return map[string]string{"id": "wrong type"}, apis.ErrWrongParamsList
			}
			id = int32(i)
		}

		err := table.SelectAndRunEach(ctx,
			func(row []interface{}, columns []dbEngine.Column) error {
				for i, col := range columns {
					name := col.Name()
					if col.BasicType() == types.UnsafePointer {
						rows[name] = renderImg(id, table.Name(), name)
					} else {
						rows[name] = row[i]
					}
				}

				return nil
			},
			dbEngine.WhereForSelect("id"),
			dbEngine.ArgsForSelect(id))
		if err != nil {
			return nil, errors.Wrap(err, "select")
		}

		return rows, nil
	}
}

// TableView show table view of table data
func TableView(preRoute string, DB *dbEngine.DB, table, patternList dbEngine.Table, priColumns []string) apis.ApiRouteHandler {
	colSelect := make([]string, 0)
	for _, col := range table.Columns() {
		if col.BasicType() == types.UnsafePointer {
			colSelect = append(colSelect, `'null' :: bytea as `+col.Name())
		} else {
			colSelect = append(colSelect, col.Name())
		}
	}

	orderBy := "order by " + strings.Join(priColumns, ",")
	conn := DB.Conn
	sql := fmt.Sprintf("select %s from %s ", strings.Join(colSelect, ","), table.Name())
	isOneID := len(priColumns) == 1 && priColumns[0] == "id"

	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		isView := bytes.Contains(ctx.Request.URI().Path(), []byte("/view"))
		user := auth.GetUserData(ctx)

		href := preRoute + table.Name() + "/form"
		// if table.Name() == "forms" {
		// 	href = preRoute + table.Name() + "/form_editor"
		// }

		href += "?html"

		lang := ctx.UserValue("lang")
		if l, ok := lang.(string); ok {
			href += "&lang=" + l
		}

		counter, _ := ctx.UserValue(ParamsCounter.Name).(int)
		firstRows := ""
		if !isView && (table.Name() == "firma") {
			for id := range user.Companies {
				if firstRows == "" {
					firstRows = " where id=ANY(ARRAY["
				} else {
					firstRows += ","
				}
				firstRows += strconv.Itoa(int(id))
			}

			if firstRows > "" {
				firstRows += "]) "
			}
		}

		firstRows += orderBy
		if counter > 0 {
			firstRows += fmt.Sprintf(" fetch first %d row only", counter)
		}

		logs.DebugLog(firstRows, user)
		rows := make([][]interface{}, 0)

		err := conn.SelectAndRunEach(ctx,
			func(row []interface{}, columns []dbEngine.Column) error {

				id := int32(0)
				vars := ""
				for i, col := range columns {
					name := col.Name()
					if isOneID && col.Name() == "id" {
						id = row[i].(int32)
					} else if colTable := table.FindColumn(name); colTable != nil && colTable.Primary() {
						vars += fmt.Sprintf(`&%s=%v`, col.Name(), row[i])
					}

					if col.BasicType() == types.UnsafePointer {
						if isView {
							row[i] = td.TdImgView(preRoute, table.Name(), name, int(id))
						} else {
							row[i] = td.TdImgBrowse(preRoute, table.Name(), name, int(id))
						}
					} else if v := GetForeignName(ctx, DB, col, row[i]); v != nil {
						foreignName := strings.TrimPrefix(name, "id_")
						if foreignName == "photos" {
							row[i] = renderImg(row[i].(int32), foreignName, "blob")
						} else {
							s := fmt.Sprintf(`<a href='%s%s/%v'>%s</a>`, preRoute,
								foreignName, row[i], v)
							row[i] = s
						}
					} else {
						row[i] = ToStandartColumnValueType(table.Name(), name, id, row[i])
					}
				}

				if !isView {
					if isOneID {
						if table.Name() == "firma" {
							c, ok := user.Companies[id]
							if ok {
								g, ok := c["edit"]
								if ok && (g == "all" || strings.Contains(g, "company")) {
									row[0] = fmt.Sprintf(`<a href='%s&id=%d'>§</a>`, href, id)
								}
							}
						} else {
							row[0] = fmt.Sprintf(`<a href='%s&id=%d' title="%[2]d">§</a>`, href, id)
						}
					} else {
						row = append([]interface{}{fmt.Sprintf(`<a href='%s%s'>§</a>`, href, vars)}, row...)
					}
				}

				rows = append(rows, row)

				return nil
			},
			sql+firstRows,
		)
		if err != nil {
			return nil, errors.Wrap(err, "select")
		}

		views.WriteHeadersHTML(ctx)

		columnDecors := CreateColumnDecors(ctx, table, patternList)
		if !isView {
			linkFormNew := fmt.Sprintf(`<a href='%s'>+</a>`, href)
			if isOneID {
				columnDecors[0].LinkNew = linkFormNew
			} else {
				colDec := forms.NewColumnDecor(dbEngine.NewStringColumn("#", "", false), patternList)
				colDec.LinkNew = linkFormNew
				columnDecors = append([]*forms.ColumnDecor{colDec}, columnDecors...)
			}
		}
		routeTable.WriteTableRow(ctx, columnDecors, rows)

		return nil, nil
	}
}

func renderImg(id int32, s ...string) string {
	name := "blob"
	if len(s) > 1 {
		name = s[1]
	}
	return fmt.Sprintf(`<img src='%s' id='img%s%v' style='max-width:100%%'/>`,
		PhotosURL(id, s...), name, id)
}

func CreateColumnDecors(ctx *fasthttp.RequestCtx, table dbEngine.Table, patternList dbEngine.Table) []*forms.ColumnDecor {
	return ToColumnDecors(ctx, table.Columns(), patternList)
}

func ToColumnDecors(ctx *fasthttp.RequestCtx, columns []dbEngine.Column, patternList dbEngine.Table) []*forms.ColumnDecor {
	colDecors := make([]*forms.ColumnDecor, len(columns))
	for i, col := range columns {
		colDecors[i] = forms.NewColumnDecor(col, patternList)
		if colDecors[i].PlaceHolder > "" {
			colDecors[i].PlaceHolder = Translate(ctx, colDecors[i].PlaceHolder)
		}

		if colDecors[i].Label > "" {
			colDecors[i].Label = Translate(ctx, colDecors[i].Label)
		} else {
			logs.DebugLog(col.Comment())
			colDecors[i].Label = Translate(ctx, col.Name())
		}
	}
	return colDecors
}

func PhotosURL(id int32, s ...string) string {
	// todo get host from ctx
	switch len(s) {
	case 0:
		return fmt.Sprintf(PathVersion+"/blob/photos?id=%d", id)
	case 1:
		return fmt.Sprintf(PathVersion+"/blob/%s?id=%d", s[0], id)
	default:
		return fmt.Sprintf(PathVersion+"/blob/%s?id=%d&name=%s", s[0], id, s[1])
	}
}

func PdfURL(id int32, s ...string) string {
	// todo get host from ctx
	switch len(s) {
	case 0:
		return fmt.Sprintf(PathVersion+"/pdf/pdf?id=%d", id)
	case 1:
		return fmt.Sprintf(PathVersion+"/pdf/%s?id=%d", s[0], id)
	default:
		return fmt.Sprintf(PathVersion+"/pdf/%s?id=%d&name=%s", s[0], id, s[1])
	}
}
