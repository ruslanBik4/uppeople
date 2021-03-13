// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type DTOContact struct {
	Id              int32         `json:"id"`
	Name            string        `json:"name"`
	Email           string        `json:"email"`
	Phone           string        `json:"phone"`
	Skype           string        `json:"skype"`
	SelectPlatforms SelectedUnits `json:"selectedPlatforms"`
	IsChecked       bool          `json:"isChecked"`
}

func (d *DTOContact) GetValue() interface{} {
	return d
}

func (d *DTOContact) NewValue() interface{} {
	return &DTOContact{}
}

func HandleAddContactForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	idCompany, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOContact)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	allPlatforms := 0
	if u.IsChecked {
		allPlatforms = 1
	}

	contacts, _ := db.NewContacts(DB)
	idC, err := contacts.Insert(ctx,
		dbEngine.ColumnsForSelect("company_id", "name", "email", "phone", "skype", "all_platforms"),
		dbEngine.ArgsForSelect(idCompany, u.Name, u.Email, u.Phone, u.Skype, allPlatforms),
	)
	if err != nil {
		return createErrResult(err)
	}

	if !u.IsChecked {
		table, _ := db.NewContacts_to_platforms(DB)
		for _, val := range u.SelectPlatforms {
			_, err := table.Insert(ctx,
				dbEngine.ColumnsForSelect("contact_id", "platform_id"),
				dbEngine.ArgsForSelect(idC, val.Id),
			)
			if err != nil {
				logs.ErrorLog(err, "NewContacts_to_platforms %s", val.Label)
			}
		}
	}

	toLogCompany(ctx, DB, idCompany, " добавил новый контакт  "+u.Name, CODE_LOG_UPDATE)

	u.Id = int32(idC)

	return u, nil
}

func HandleEditContactForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	idCompany, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOContact)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	allPlatforms := 0
	if u.IsChecked {
		allPlatforms = 1
	}

	contacts, _ := db.NewContacts(DB)
	idC, err := contacts.Update(ctx,
		dbEngine.ColumnsForSelect("company_id", "name", "email", "phone", "skype", "all_platforms"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(idCompany, u.Name, u.Email, u.Phone, u.Skype, allPlatforms, u.Id),
	)
	if err != nil {
		return createErrResult(err)
	}

	if !u.IsChecked {
		table, _ := db.NewContacts_to_platforms(DB)
		for _, val := range u.SelectPlatforms {
			_, err := table.Upsert(ctx,
				dbEngine.ColumnsForSelect("contact_id", "platform_id"),
				dbEngine.ArgsForSelect(u.Id, val.Id),
			)
			if err != nil {
				logs.ErrorLog(err, "NewContacts_to_platforms %s", val.Label)
			}
		}
	}

	toLogCompany(ctx, DB, idCompany, " edit контакт  "+u.Name, CODE_LOG_UPDATE)

	u.Id = int32(idC)

	return u, nil
}

func HandleDeleteContactForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	err := DB.Conn.ExecDDL(ctx, "delete from contacts where id=$1", id)
	if err != nil {
		return createErrResult(err)
	}

	toLogCompany(ctx, DB, id, " удалил контакт  ", CODE_LOG_DELETE)

	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}

type ViewContact struct {
	*db.ContactsFields
	SelectPlatforms SelectedUnits `json:"selectedPlatforms"`
}

func HandleViewContactForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: "wrong type, expect int32",
		}, apis.ErrWrongParamsList
	}

	contacts, _ := db.NewContacts(DB)
	err := contacts.SelectOneAndScan(ctx,
		contacts,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	v := &ViewContact{
		contacts.Record,
		SelectedUnits{},
	}
	err = DB.Conn.SelectAndScanEach(ctx,
		nil,
		&v.SelectPlatforms,
		`select p.id, p.nazva as label, p.nazva as value 
			from contacts_to_platforms c join platforms p on p.id=platform_id
			where contact_id=$1`,
		id,
	)
	if err != nil {
		return createErrResult(err)
	}

	return v, nil
}
