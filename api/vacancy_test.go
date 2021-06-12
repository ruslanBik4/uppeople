// Copyright 2020 Author: Ruslan Bikchentaev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"

	. "github.com/ruslanBik4/uppeople/db"
)

func TestVacancyDTO_GetValue(t *testing.T) {
	type fields struct {
		VacanciesFields       *VacanciesFields
		Comment               string
		Description           string
		Phone                 string
		Status                string
		SelectCompany         SelectedUnit
		SelectLocation        SelectedUnit
		SelectPlatform        SelectedUnit
		SelectSeniority       SelectedUnit
		SelectRecruiter       []SelectedUnit
		SelectedVacancyStatus int32
	}
	var tests = []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
		{
			`{"selectPlatform":{"id":1,"label":"Java","value":"java"},"selectSeniority":{"id":2,"label":"Mid","value":"mid"},"selectCompany":{"id":2,"label":"Voicespin","value":"voicespin"},"selectLocation":{"id":2,"label":"Kyiv","value":"kyiv"},"selectRecruiter":[{"id":12,"label":"Helga Nizhnyk","value":"helga nizhnyk"},{"id":11,"label":"Ed","value":"ed"},{"id":19,"label":"Kateryna Denysenko","value":"kateryna denysenko"}],"salary":1000,"comment":"","link":"","selectedVacancyStatus":1,"description":"<p>test</p>\n","details":"<p>test</p>\n"}`,
			fields{},
			&VacancyDTO{
				SelectPlatform:  SelectedUnit{1, "Java", "java"},
				SelectSeniority: SelectedUnit{2, "Mid", "mid"},
				SelectCompany:   SelectedUnit{2, "Voicespin", "voicespin"},
				SelectLocation:  SelectedUnit{2, "Kyiv", "kyiv"},
				SelectRecruiter: []SelectedUnit{
					{
						12,
						"Helga Nizhnyk",
						"helga nizhnyk",
					},
					{11, "Ed", "ed"},
					{19, "Kateryna Denysenko", "kateryna denysenko"},
				},
				SelectedVacancyStatus: 1,
				VacanciesFields: &VacanciesFields{
					Salary:      1000,
					Description: "<p>test</p>\n",
					Details:     "<p>test</p>\n",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VacancyDTO{
				VacanciesFields:       tt.fields.VacanciesFields,
				Comment:               tt.fields.Comment,
				SelectCompany:         tt.fields.SelectCompany,
				SelectLocation:        tt.fields.SelectLocation,
				SelectPlatform:        tt.fields.SelectPlatform,
				SelectSeniority:       tt.fields.SelectSeniority,
				SelectRecruiter:       tt.fields.SelectRecruiter,
				SelectedVacancyStatus: tt.fields.SelectedVacancyStatus,
			}
			err := jsoniter.UnmarshalFromString(tt.name, &v)
			if err != nil {
				t.Error(err, "")
			}
			if got := v.GetValue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
