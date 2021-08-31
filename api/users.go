// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"fmt"

	"github.com/ruslanBik4/dbEngine/dbEngine"
	"github.com/ruslanBik4/httpgo/apis"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"

	"github.com/ruslanBik4/uppeople/db"
)

type UserResponse struct {
	Users       UserRows         `json:"users"`
	Partners    db.SelectedUnits `json:"partners"`
	Freelancers db.SelectedUnits `json:"freelancers"`
	Recruiters  db.SelectedUnits `json:"recruiters"`
}

type DTOUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     int32  `json:"role"`
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

func NewUserRow() *UserRow {
	return &UserRow{&db.UsersFields{}, 0, 0, 0}
}

func (u *UserRow) GetFields(columns []dbEngine.Column) []interface{} {
	row := make([]interface{}, len(columns))
	for i, col := range columns {
		switch col.Name() {
		case "create_count":
			row[i] = &u.CreateCount
		case "update_count":
			row[i] = &u.UpdateCount
		case "send_count":
			row[i] = &u.SendCount
		default:
			row[i] = u.RefColValue(col.Name())
		}
	}

	return row
}

type UserRows []*UserRow

func (a *UserRows) GetFields(columns []dbEngine.Column) []interface{} {
	row := NewUserRow()
	*a = append(*a, row)

	return row.GetFields(columns)
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

	return users.Record, nil
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

	columns := []string{"name", "email", "phone", "role_id"}
	args := []interface{}{u.Name, u.Email, u.Phone, u.Role, u.Id}

	if u.Password > "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		args = append(args, hash)
		columns = append(columns, "password")
	}

	users, _ := db.NewUsers(DB)
	i, err := users.Update(ctx,
		dbEngine.ColumnsForSelect(columns...),
		dbEngine.WhereForSelect("id"),
		dbEngine.ArgsForSelect(args...),
	)
	if err != nil {
		return createErrResult(err)
	}

	return createResult(i)
}

func HandleDelUser(ctx *fasthttp.RequestCtx) (interface{}, error) {
	id, ok := ctx.UserValue(ParamID.Name).(int32)
	if !ok {
		return map[string]string{
			ParamID.Name: fmt.Sprintf("wrong type %T, expect int32 ", ctx.UserValue(ParamID.Name)),
		}, apis.ErrWrongParamsList
	}

	DB, ok := ctx.UserValue("DB").(*dbEngine.DB)
	if !ok {
		return nil, dbEngine.ErrDBNotFound
	}

	err := DB.Conn.ExecDDL(ctx, "delete from users where id = $1", id)
	if err != nil {
		return createErrResult(err)
	}

	ctx.SetStatusCode(fasthttp.StatusAccepted)

	return nil, nil
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

	r := make(UserRows, 0)
	users, _ := DB.Tables["all_staff"]
	err := users.SelectAndScanEach(ctx,
		nil,
		&r,
		dbEngine.OrderBy("name"),
	)
	if err != nil {
		return createErrResult(err)
	}

	return UserResponse{
		Users: r,
	}, nil
}
