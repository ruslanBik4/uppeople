// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/pkg/errors"
	"github.com/ruslanBik4/dbEngine/dbEngine"
	_go "github.com/ruslanBik4/dbEngine/generators/go"
	"github.com/ruslanBik4/logs"

	"github.com/ruslanBik4/uppeople/db"
)

func main() {

	DB := db.GetDB(nil)
	if DB == nil {
		logs.Fatal(dbEngine.ErrDBNotFound, " ot DB")
	}

	creator, err := _go.NewCreator("./db")
	if err != nil {
		logs.ErrorLog(errors.Wrap(err, "NewCreator"))
		return
	}

	for name, table := range DB.Tables {
		err = creator.MakeStruct(table)
		if err != nil {
			logs.ErrorLog(errors.Wrap(err, "makeStruct - "+name))
		}
	}

}
