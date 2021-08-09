// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"database/sql"
	"math"
	"math/big"
	"strings"
	"unsafe"

	"github.com/jackc/pgtype"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/ruslanBik4/logs"
)

func ReadNullString(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	val := (*sql.NullString)(ptr)

	switch t := iter.WhatIsNext(); t {
	case jsoniter.NilValue:
		val.Valid = !iter.ReadNil()
		if val.Valid {
			val.String = "wrong value, nil is not NULL!"
		}

	case jsoniter.StringValue:
		src := iter.ReadString()
		err := val.Scan(src)
		if err != nil {
			logs.ErrorLog(err, val, src)
		}

	case jsoniter.ArrayValue:
		val.Valid = iter.ReadArrayCB(func(iterator *jsoniter.Iterator) bool {
			val.String = iter.ReadString()
			return true
		})

	case jsoniter.ObjectValue:
		val.Valid = iter.ReadMapCB(func(iterator *jsoniter.Iterator, key string) bool {
			switch strings.ToLower(key) {
			case "string":
				val.String = iter.ReadString()
				return true
			case "valid":
				val.Valid = iter.ReadBool()
				return val.Valid
			default:
				logs.ErrorLog(errors.New("unknown key of NUllString"), key)
				return false
			}
		})

	default:
		logs.ErrorLog(errors.New("unknown type"), t)
	}
}

func init() {
	jsoniter.RegisterTypeEncoderFunc("pgtype.Numrange",
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			val := (*pgtype.Numrange)(ptr)
			if val.Status == pgtype.Present {
				stream.WriteArrayStart()
				if val.Lower.Int != nil {
					stream.WriteFloat64Lossy(float64(val.Lower.Int.Int64()) * math.Pow10(int(val.Lower.Exp)))
				} else if val.Lower.NaN {
					stream.WriteFloat64Lossy(-math.MaxFloat64)
				} else {
					stream.WriteFloat64Lossy(0)
				}

				stream.WriteMore()
				if val.Upper.Int != nil {
					stream.WriteFloat64Lossy(float64(val.Upper.Int.Int64()) * math.Pow10(int(val.Upper.Exp)))
				} else if val.Upper.NaN {
					stream.WriteFloat64Lossy(math.MaxFloat64)
				} else {
					stream.WriteFloat64Lossy(0)
				}

				stream.WriteArrayEnd()
			} else {
				stream.WriteEmptyArray()
			}
		},
		func(pointer unsafe.Pointer) bool {
			return false
		})

	jsoniter.RegisterTypeDecoderFunc("pgtype.Numrange",
		func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
			i := 0
			numrange := (*pgtype.Numrange)(ptr)
			numrange.Status = pgtype.Undefined
			var err error
			iter.ReadArrayCB(func(iter *jsoniter.Iterator) bool {
				val := iter.ReadFloat64()
				f := big.NewFloat(val)
				if i == 0 {
					err = numrange.Lower.Scan(f.String())
				} else {
					err = numrange.Upper.Scan(f.String())
				}

				i++
				if err != nil {
					logs.ErrorLog(err, val)
				} else if i > 1 {
					numrange.Status = pgtype.Present
					numrange.UpperType = pgtype.Inclusive
					numrange.LowerType = pgtype.Inclusive
				}

				return true
			})

		})

	jsoniter.RegisterTypeDecoderFunc("sql.NullString", ReadNullString)
	jsoniter.RegisterTypeDecoderFunc("*sql.NullString", ReadNullString)
}
