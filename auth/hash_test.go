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
		key  string
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
				"$2a$10$fHY1cLCnHEWW4wogO3bJE.kwOdPZst.kyIUEhJKx9KqCbuJi7Vlza",
				"alona.ryd",
			},
			false,
		},
		{
			"Anastasiya Syradoeva",
			args{
				"$2y$10$phuDw8Vd6AtaDeKVkmW/4.KWsShN6Nwi2H/eMcpcfv0AP.UQpWBIS",
				"syradoevanast",
			},
			false,
		},
		{
			"Leonova Kate",
			args{
				"$2a$10$ivzIGJ6g7jhUiN/eWpVfpe82sbj8pIrQ/u/nbyRDOwDlUw.tcfW22",
				"leonovakate",
			},
			false,
		},
		{
			"Hanna Skorokhod (макбук)",
			args{
				"$2a$10$X9OPL0.x9FGqUTO53d4tou.CFKF1jbqNgVpXLb7j4DW8QDi41JlC.",
				"hanna.skorokhod",
			},
			false,
		},
		{
			"mironyana90@gmail.com",
			args{
				"$2a$10$X9OPL0.x9FGqUTO53d4tou.CFKF1jbqNgVpXLb7j4DW8QDi41JlC.",
				"hanna.skorokhod",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !assert.Equal(t, tt.wantErr, CheckPass(tt.args.hash, tt.args.key) != nil) {
				h, err := NewHash(tt.args.key)
				if assert.Nil(t, err) {
					// b, err := bcrypt.GenerateFromPassword([]byte(tt.args.pass), bcrypt.DefaultCost)
					// assert.Nil(t, err)
					assert.Equal(t, tt.args.hash, string(h))
				}

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
