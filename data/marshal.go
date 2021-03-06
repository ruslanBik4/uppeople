// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"database/sql"
	"math"
	"math/big"
	"unsafe"

	"github.com/jackc/pgtype"
	jsoniter "github.com/json-iterator/go"
	"github.com/ruslanBik4/logs"
)

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

	jsoniter.RegisterTypeEncoderFunc("sql.NullString",
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			val := (*sql.NullString)(ptr)
			if val.Valid {
				stream.WriteString(val.String)
			} else {
				stream.WriteString("nil")
			}
		},
		func(pointer unsafe.Pointer) bool {
			return false
		})

	jsoniter.RegisterTypeEncoderFunc("sql.NullInt32",
		func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
			val := (*sql.NullInt32)(ptr)
			if val.Valid {
				stream.WriteInt32(val.Int32)
			} else {
				stream.WriteString("nil")
			}
		},
		func(pointer unsafe.Pointer) bool {
			return false
		})
}
