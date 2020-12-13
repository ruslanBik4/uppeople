// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"fmt"
	"path"
	"strings"

	"github.com/ruslanBik4/httpgo/views"
	"golang.org/x/net/context"

	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/httpgo/views/templates/forms"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type FormsList struct {
	listForms   map[string]string
	DB          *dbEngine.DB
	patternList *db.Patterns_list
	customForms *db.Forms
	formBlocks  *db.Form_blocks
	formFields  *db.Form_fields
	formEditor  *forms.FormField
	blockForms  forms.BlockColumns
	blockBlock  forms.BlockColumns
	blockField  forms.BlockColumns
	Routes      apis.ApiRoutes
}

func NewFormsList() *FormsList {
	return &FormsList{listForms: make(map[string]string, 0)}
}

func (fl *FormsList) Add(path, name string) {
	fl.listForms[path] = name
}

func (fl *FormsList) AddCustomForms(DB *dbEngine.DB, routes apis.ApiRoutes, pathVersion string, patternList *db.Patterns_list) error {
	customForms, err := db.NewForms(DB)
	if err != nil {
		return errors.Wrap(err, "db.NewForms")
	}

	fl.customForms = customForms
	fl.formBlocks, err = db.NewForm_blocks(DB)
	if err != nil {
		return errors.Wrap(err, "db.NewForms")
	}

	fl.formFields, err = db.NewForm_fields(DB)
	if err != nil {
		return errors.Wrap(err, "db.NewForms")
	}

	fl.DB = DB
	fl.patternList = patternList

	ctx := context.Background()
	const nameEditor = "form_editor"
	err = customForms.SelectSelfScanEach(ctx, func(record *db.FormsFields) error {
		form := fl.getFormFields(record)

		pathForm := path.Join(pathVersion, "forms", record.Title)
		if record.Title == nameEditor {
			routes[pathForm] = &apis.ApiRoute{
				Desc: nameEditor,
				Fnc:  fl.HandleFormEditor(form),
				Params: []apis.InParam{
					paramsID,
					ParamsLang,
					ParamsHTML,
				},
			}
		} else {
			routes[pathForm] = &apis.ApiRoute{
				Desc: record.Description,
				Fnc:  fl.HandleShowCustomForm(DB, form),
				Params: []apis.InParam{
					paramsID,
					ParamsLang,
					ParamsHTML,
				},
			}
		}

		fl.Add(pathForm, record.Title)

		return nil
	})

	return err
}

func (fl *FormsList) getFormFields(record *db.FormsFields) forms.FormField {
	form := forms.FormField{
		Title:       record.Title,
		Action:      record.Action,
		Method:      "POST",
		Description: record.Description,
		HideBlock:   record.HideBlock,
	}

	if !record.Post {
		form.Method = "GET"
	}

	return form
}

func (fl *FormsList) HandleFormEditor(form forms.FormField) apis.ApiRouteHandler {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {

		err := fl.customForms.SelectSelfScanEach(ctx, func(record *db.FormsFields) error {
			_, blocks, err := fl.getBlocks(ctx, record.Id)
			if err != nil {
				return err
			}

			fl.formEditor = &form
			for _, block := range blocks {
				switch block.Title {
				case "Form setting":
					fl.blockForms = block
				case "Block setting":
					fl.blockBlock = block
				case "Fields setting":
					fl.blockField = block
				}
			}

			blockField := []forms.BlockColumns{
				fl.blockForms, fl.blockBlock,
			}

			for _, col := range fl.blockForms.Columns {
				col.Value = record.ColValue(col.Name())
			}

			err = fl.formBlocks.SelectSelfScanEach(ctx,
				func(record *db.Form_blocksFields) error {

					for _, col := range fl.blockBlock.Columns {
						col.Value = record.ColValue(col.Name())
					}

					err := fl.formFields.SelectSelfScanEach(ctx,
						func(record *db.Form_fieldsFields) error {

							block := forms.BlockColumns{
								Buttons:     nil,
								Columns:     make([]*forms.ColumnDecor, len(fl.blockField.Columns)),
								Id:          0,
								Multiple:    false,
								Title:       "",
								Description: "",
							}
							for i, col := range fl.blockField.Columns {
								block.Columns[i] = copyColumnDecor(ctx, col)
								block.Columns[i].Value = record.ColValue(col.Name())
							}

							blockField = append(blockField, block)

							return nil
						},
						dbEngine.WhereForSelect("id_form_blocks"),
						dbEngine.ArgsForSelect(record.Id))
					if err != nil {
						return errors.Wrap(err, "form fields getting")
					}

					return nil
				}, dbEngine.WhereForSelect("id_forms"), dbEngine.ArgsForSelect(record.Id))
			if err != nil {
				return errors.Wrap(err, "form block getting")
			}

			_, ok := ctx.UserValue("html").(bool)
			if !ok {
				views.WriteJSONHeaders(ctx)
			}

			fl.formEditor.WriteRenderForm(
				ctx.Response.BodyWriter(),
				ok,
				blockField...,
			)

			return nil
		},
			dbEngine.WhereForSelect("title"),
			dbEngine.ArgsForSelect(form.Title),
		)
		if err != nil {
			return nil, errors.Wrap(err, "SelectSelfScanEach")
		}

		return nil, nil
	}
}

func (fl *FormsList) getBlocks(ctx *fasthttp.RequestCtx, id int32) ([]dbEngine.Table, []forms.BlockColumns, error) {

	tables := make([]dbEngine.Table, 0)
	blocks := make([]forms.BlockColumns, 0)

	err := fl.formBlocks.SelectSelfScanEach(ctx, func(recordBlock *db.Form_blocksFields) error {
		block := forms.BlockColumns{
			// Buttons:     recordBlock.Buttons,
			Columns:     make([]*forms.ColumnDecor, 0),
			Id:          int(recordBlock.Id),
			Multiple:    true,
			Title:       recordBlock.Title,
			Description: recordBlock.Description,
		}

		table, ok := fl.DB.Tables[recordBlock.Tablename]
		if !ok {
			return dbEngine.ErrNotFoundTable{
				Table: recordBlock.Tablename,
			}
		}

		btns, ok := recordBlock.Buttons.([]interface{})
		if ok {
			for _, part := range btns {
				btn, ok := part.(map[string]interface{})
				if ok {
					block.Buttons = append(block.Buttons,
						forms.Button{
							Title:      Translate(ctx, btn["title"].(string)),
							Position:   false,
							ButtonType: "submit",
						},
					)
				} else {
					logs.DebugLog("%T %+[1]v", btn)
				}

			}
		} else if recordBlock.Buttons != nil {
			logs.DebugLog("%T %+[1]v", recordBlock.Buttons)
		}
		tables = append(tables, table)

		if len(recordBlock.Columns) > 0 {
			block.Columns = fl.getColumns(table, recordBlock.Columns)
		}

		err := fl.formFields.SelectSelfScanEach(ctx,
			func(record *db.Form_fieldsFields) error {
				col := table.FindColumn(record.Name)
				if col == nil {
					logs.DebugLog("%s haven't field '%s'. I created string col instead", table.Name(), record.Name)
					col = dbEngine.NewStringColumn(record.Name, record.Title, record.Required)
				}

				colDec := forms.NewColumnDecor(col, fl.patternList)
				colDec.IsReadOnly = record.Readonly
				colDec.Label = Translate(ctx, record.Title)
				colDec.PlaceHolder = Translate(ctx, record.Placeholder)
				// colDec.IsSlice = true
				if record.Input_type == "select" {
					switch col.Name() {
					case "tablename":
						colDec.SelectOptions = make(map[string]string)
						for name, table := range fl.DB.Tables {
							comment := table.Comment()
							if strings.TrimSpace(comment) == "" {
								comment = table.Name()
							}

							colDec.SelectOptions[comment] = name
						}
					case "name":
						colDec.SelectOptions = make(map[string]string)
						for _, col := range table.Columns() {
							comment := col.Comment()
							if strings.TrimSpace(comment) == "" {
								comment = col.Name()
							} else if p := strings.Split(comment, "{"); len(p) > 1 {
								comment = p[0]
							}

							colDec.SelectOptions[comment] = col.Name()
						}
					case "pattern":
						colDec.SelectOptions = make(map[string]string)
						err := fl.patternList.SelectSelfScanEach(context.Background(),
							func(record *db.Patterns_listFields) error {
								colDec.SelectOptions[record.Name] = record.Name
								return nil
							},
						)
						if err != nil {
							logs.ErrorLog(err, "")
						}

					case "input_type":
						colDec.SelectOptions = map[string]string{
							"Dropdown list of titles": "select",
							"Simple text":             "text",
							"Long text":               "textarea",
							"password":                "password",
							"datetime":                "datetime",
							"number":                  "number",
							"Simple search":           "search",
							"email":                   "email",
							"Phone":                   "tel",
							"Link URL":                "url",
							"Switch":                  "switch",
						}
					default:
						GetForeignOptions(ctx, fl.DB, colDec, nil)
					}
				}

				block.Columns = append(block.Columns, colDec)

				return nil

			},
			dbEngine.WhereForSelect("id_form_blocks"),
			dbEngine.ArgsForSelect(recordBlock.Id),
			dbEngine.OrderBy("id"),
		)
		if err != nil {
			return errors.Wrap(err, "form fields getting")
		}

		blocks = append(blocks, block)
		return nil
	},
		dbEngine.WhereForSelect("id_forms"),
		dbEngine.ArgsForSelect(id),
		dbEngine.OrderBy("id"),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "form block getting")
	}

	return tables, blocks, nil
}

func (fl *FormsList) getColumns(table dbEngine.Table, columns []string) []*forms.ColumnDecor {
	colDecors := make([]*forms.ColumnDecor, 0)
	for _, name := range columns {
		col := table.FindColumn(name)
		if col == nil {
			continue
		}
		if !col.Primary() && !strings.Contains(col.Comment(), " (read_only)") && (col.Name() != "id") {
			colDec := forms.NewColumnDecor(col, fl.patternList)
			colDec.Value = col.Default()
			colDec.IsSlice = (strings.HasPrefix(col.Type(), "_"))

			if colDec.Value == "NULL" {
				colDec.Value = nil
			}

			if col.Type() == "text" {
				colDec.InputType = "textarea"
			}

			colDecors = append(colDecors, colDec)
		}
	}
	return colDecors
}

func (fl *FormsList) HandleShowCustomForm(DB *dbEngine.DB, form forms.FormField) apis.ApiRouteHandler {
	return func(ctx *fasthttp.RequestCtx) (interface{}, error) {
		// lang := ctx.UserValue(ParamsLang.Name)
		id, ok := ctx.UserValue("id").(int32)
		err := fl.customForms.SelectSelfScanEach(ctx, func(record *db.FormsFields) error {
			form = fl.getFormFields(record)

			tables, blocks, err := fl.getBlocks(ctx, record.Id)
			if err != nil {
				return errors.Wrap(err, "getBlocks")
			}

			for key, table := range tables {
				if ok {
					err := table.SelectAndRunEach(ctx, func(values []interface{}, columns []dbEngine.Column) error {
						ok = false
						for _, colDecor := range blocks[key].Columns {
							name := colDecor.Name()
							for i, col := range columns {
								if name == col.Name() {
									colDecor.Value = values[i]
									GetForeignOptions(ctx, DB, colDecor, id)
									if name == "id" {
										colDecor.IsHidden = true
									} else if col.Type() == "string" {
										val, isStr := values[i].(string)
										if isStr && (val > "") {
											colDecor.Value = Translate(ctx, val)
										} else if isStr {
											colDecor.Value = ""
										}
									}
								}
							}
						}

						return nil
					}, dbEngine.WhereForSelect("id"), dbEngine.ArgsForSelect(id))
					if err != nil {
						return err
					}

					if ok {
						ctx.SetStatusCode(fasthttp.StatusNoContent)
						return nil
					}

				} else {
					for _, col := range blocks[key].Columns {
						GetForeignOptions(ctx, DB, col, nil)
					}
				}
			}

			btnList := []forms.Button{
				{ButtonType: "submit", Title: "Insert", Position: true},
				{ButtonType: "reset", Title: "Clear", Position: false},
			}

			blocks = append(blocks, forms.BlockColumns{
				Buttons:     btnList,
				Id:          -1,
				Title:       "Apply changes",
				Description: "",
			})

			_, ok := ctx.UserValue("html").(bool)
			if !ok {
				views.WriteJSONHeaders(ctx)
			}

			form.WriteRenderForm(
				ctx.Response.BodyWriter(),
				ok,
				blocks...)
			return err

		},
			dbEngine.WhereForSelect("title"),
			dbEngine.ArgsForSelect(form.Title),
		)
		if err != nil {
			return nil, errors.Wrap(err, "SelectSelfScanEach")
		}

		return nil, nil
	}
}

func (fl *FormsList) GetRoute() *apis.ApiRoute {
	return &apis.ApiRoute{
		Desc: "db routes list",
		Fnc: func(ctx *fasthttp.RequestCtx) (interface{}, error) {
			// todo make json
			s, comma := "[", ""
			for key, val := range fl.listForms {
				s += fmt.Sprintf(comma+`{"title":"%s", "value":"%s"}`, strings.Title(val), key)
				comma = ","
			}

			return s + "]", nil
		},
	}
}
