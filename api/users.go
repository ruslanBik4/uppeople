// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/ruslanBik4/logs"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruslanBik4/uppeople/db"
)

type UserResponse struct {
	Users       []UserRow     `json:"users"`
	Partners    SelectedUnits `json:"partners"`
	Freelancers SelectedUnits `json:"freelancers"`
	Recruiters  SelectedUnits `json:"recruiters"`
}

type DTOUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

func (d *DTOUser) GetValue() interface{} {
	return d
}

func (d *DTOUser) NewValue() interface{} {
	return &DTOUser{}
}

type UserRow struct {
	*db.UsersFields
	CreateCount int32 `json:"createCount"`
	UpdateCount int32 `json:"updateCount"`
	SendCount   int32 `json:"sendCount"`
}

func HandleGetUser(ctx *fasthttp.RequestCtx) (interface{}, error) {
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
	users, _ := db.NewUsers(DB)
	err := users.SelectOneAndScan(ctx,
		users,
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(id),
	)
	if err != nil {
		return createErrResult(err)
	}

	return users.Table, nil
}

func HandleEditUser(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOUser)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	users, _ := db.NewUsers(DB)
	i, err := users.Update(ctx,
		dbEngine.ColumnsForSelect("name", "email", "phone", "role_id"),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(u.Name, u.Email, u.Phone, u.Role, u.Id),
	)
	if err != nil {
		return createErrResult(err)
	}

	return createResult(i)
}

func HandleNewUser(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	u, ok := ctx.UserValue(apis.JSONParams).(*DTOUser)
	if !ok {
		return "wrong DTO", apis.ErrWrongParamsList
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	users, _ := db.NewUsers(DB)
	id, err := users.Insert(ctx,
		dbEngine.ColumnsForSelect("name", "email", "phone", "role_id", "password"),
		dbEngine.ArgsForSelect(u.Name, u.Email, u.Phone, u.Role, hash),
	)
	if err != nil {
		return createErrResult(err)
	}

	ctx.SetStatusCode(fasthttp.StatusCreated)

	return createResult(id)
}

func HandleAllStaff(ctx *fasthttp.RequestCtx) (interface{}, error) {
	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	r := make([]UserRow, 0)
	users, _ := db.NewUsers(DB)
	err := users.SelectSelfScanEach(ctx,
		func(record *db.UsersFields) error {
			row := UserRow{users.Record, 0, 0, 0}

			err := DB.Conn.SelectOneAndScan(ctx,
				[]interface{}{&row.CreateCount, &row.UpdateCount, &row.SendCount},
				`select count(*) FILTER ( WHERE kod_deystviya = $1 ),
       count(*) FILTER ( WHERE kod_deystviya = $2 ),
       count(*) FILTER ( WHERE kod_deystviya = $3 )
from logs
where age(date_create) < interval '7 day' and user_id = $4`,
				CODE_LOG_INSERT,
				CODE_LOG_UPDATE,
				CODE_LOG_PEFORM,
				users.Record.Id,
			)
			if err != nil {
				logs.ErrorLog(err, "DB.Conn.SelectOneAndScan")
			}
			r = append(r, row)
			return nil
		},
		dbEngine.OrderBy("name"),
	)
	if err != nil {
		return createErrResult(err)
	}

	return UserResponse{
		Users: r,
	}, nil
}
