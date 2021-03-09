// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPass(hash, pass string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

func NewHash(pass string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
}
func NewHashWithCost(pass string, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pass), cost)
}
