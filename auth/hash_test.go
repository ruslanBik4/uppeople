// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCheckPass(t *testing.T) {
	type args struct {
		hash string
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"alyona",
			args{
				"$2y$10$afhoGJsACLqisT6qOzQoTOY0A4RFbHVebVenA.xIrEtEbWrrntHh.",
				"alona.ryd",
			},
			false,
		},
		{
			"Anastasiya Syradoeva",
			args{
				"$2a$10$AdHqywHa6BXWnwVHmMOZoeqaWfuxY4wYeLM7YsAYfMn9xk6yliLze",
				"syradoevanast",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !assert.Equal(t, tt.wantErr, CheckPass(tt.args.hash, tt.args.pass)) {
				h, err := NewHash(tt.args.pass)
				assert.Nil(t, err)
				t.Logf("Hash = %v", h)
			}
		})
	}
}

func TestNewHash(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"alyona",
			args{
				"alona.ryd",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewHash(tt.args.pass)
			assert.Nil(t, err)
			t.Logf("Hash = %s", h)
		})
	}
}

func TestNewHashWithCost(t *testing.T) {
	type args struct {
		pass string
		cost []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"alyona",
			args{
				"alona.ryd",
				[]int{
					bcrypt.MinCost,
					bcrypt.DefaultCost,
					bcrypt.MaxCost,
					2,
					3,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, cost := range tt.args.cost {

				h, err := NewHashWithCost(tt.args.pass, cost)
				assert.Nil(t, err)
				t.Logf("%d. Hash = %s (cost=%d)", i, h, cost)
			}
		})
	}

}
