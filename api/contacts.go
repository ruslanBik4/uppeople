// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"

	"github.com/ruslanBik4/uppeople/db"
)

type ViewContact struct {
	*db.ContactsFields
	SelectPlatforms db.SelectedUnits `json:"selectedPlatforms"`
}

type DTOContact struct {
	Id              int32            `json:"id"`
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	Phone           string           `json:"phone"`
	Skype           string           `json:"skype"`
	SelectPlatforms db.SelectedUnits `json:"selectedPlatforms"`
	IsChecked       bool             `json:"isChecked"`
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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOContact)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	var platforms []int32
	if u.IsChecked {
		for _, val := range u.SelectPlatforms {
			platforms = append(platforms, val.Id)
		}
	}

	contacts, _ := db.NewContacts(DB)
	idC, err := contacts.Insert(ctx,
		dbEngine.ColumnsForSelect("company_id", "name", "email", "phone", "skype", "platforms"),
		dbEngine.ArgsForSelect(idCompany, u.Name, u.Email, u.Phone, u.Skype, platforms),
	)
	if err != nil {
		return createErrResult(err)
	}

	if idC > 0 {
		toLogCompanyUpdate(ctx, idCompany, map[string]interface{}{"contact_id": idC, "platforms": platforms})
	}
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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOContact)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	var platforms []int32
	if u.IsChecked {
		for _, val := range u.SelectPlatforms {
			platforms = append(platforms, val.Id)
		}
	}

	contacts, _ := db.NewContacts(DB)
	i, err := contacts.Update(ctx,
		dbEngine.ColumnsForSelect("company_id", "name", "email", "phone", "skype", "platforms"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(idCompany, u.Name, u.Email, u.Phone, u.Skype, platforms, u.Id),
	)
	if err != nil {
		return createErrResult(err)
	}

	if i > 0 {
		toLogCompanyUpdate(ctx, idCompany, map[string]interface{}{"contact_id": u.Id, "platforms": platforms})
	}

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
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	err := DB.Conn.ExecDDL(ctx, "delete from contacts where id=$1", id)
	if err != nil {
		return createErrResult(err)
	}

	toLogCompanyDelete(ctx, id, fmt.Sprintf(" удалил контакт %d", id))

	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
}

func HandleViewContactForCompany(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
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
		db.SelectedUnits{},
	}
	err = DB.Conn.SelectAndScanEach(ctx,
		nil,
		&v.SelectPlatforms,
		`select p.id, p.name as label, p.name as value 
			from public.platforms p where p.id=ANY($1)`,
		contacts.Record.Platforms,
	)
	if err != nil {
		return createErrResult(err)
	}

	return v, nil
}
