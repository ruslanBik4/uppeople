// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func Test_createErrResult(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"1",
			args{
				errors.New("Key (linkedin)=(https://www.linkedin.com/in/vladislav-yena/) already exists."),
			},
			map[string]string{
				"linkedin": "`https://www.linkedin.com/in/vladislav-yena/` already exists",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createErrResult(tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("createErrResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createErrResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}
