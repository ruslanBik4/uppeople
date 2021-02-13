// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"database/sql"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func TestReadNullString(t *testing.T) {
	type testRes struct {
		V sql.NullString
	}
	tests := []struct {
		name string
		data []byte
		res  testRes
	}{
		// TODO: Add test cases.
		{
			"simple",
			[]byte(`{"v":"test"}`),
			testRes{
				sql.NullString{
					String: "test",
					Valid:  true,
				},
			},
		},
		{
			"composite",
			[]byte(`{"v":{"String": "test"}}`),
			testRes{
				sql.NullString{
					String: "test",
					Valid:  true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var res testRes
			err := jsoniter.Unmarshal(tt.data, &res)
			assert.Nil(t, err, "Must be not error")
			assert.Equal(t, tt.res.V, res.V)
		})
	}
}
