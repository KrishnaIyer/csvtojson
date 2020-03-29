// Copyright © 2020 Krishna Iyer Easwaran
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csv

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validCSV = `first_name,last_name,username
"Rob","Pike","rob"
Ken,Thompson,ken
"Robert","Griesemer","gri"`
	invalidCSV = `first_name,last_name,username
"Rob","Pike"
Ken,Thompson,ken
"Robert","Griesemer","gri"`
	validCSVWithEmptyValues = `first_name,last_name,username
"Rob",,"Pike"
Ken,Thompson,ken
"Robert","Griesemer","gri"`
)

func TestCSV(t *testing.T) {
	ctx := context.Background()
	for _, tc := range []struct {
		Name           string
		Input          string
		Marshaled      string
		AllowMalformed bool
		FillEmptyWith  string
		CheckValues    func([]map[string]string) bool
		ExpectedError  bool
	}{
		{
			Name:      "Valid",
			Input:     validCSV,
			Marshaled: `[{"first_name":"Rob","last_name":"Pike","username":"rob"},{"first_name":"Ken","last_name":"Thompson","username":"ken"},{"first_name":"Robert","last_name":"Griesemer","username":"gri"}]`,
		},
		{
			Name:          "Invalid",
			Input:         invalidCSV,
			ExpectedError: true,
		},
		{
			Name:           "AllowMalformed",
			Input:          invalidCSV,
			AllowMalformed: true,
			FillEmptyWith:  "test",
			CheckValues: func(val []map[string]string) bool {
				return val[0]["username"] == "test"
			},
		},
		{
			Name:  "ValidWithEmptyValues",
			Input: validCSVWithEmptyValues,
			CheckValues: func(val []map[string]string) bool {
				return val[0]["last_name"] == ""
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			csv, err := New(ctx, []byte(tc.Input), Config{
				AllowMalformed: tc.AllowMalformed,
				FillEmptyWith:  tc.FillEmptyWith,
			})
			if !tc.ExpectedError {
				if !assert.Nil(t, err) {
					t.Fatalf("Unexpected error: %v", err)
				}
			} else {
				if !assert.NotNil(t, err) {
					t.Fatal("Expected error but none received")
				}
			}
			if tc.CheckValues != nil {
				if !tc.CheckValues(csv.Values()) {
					t.Fatalf("Unexpected csv : %v", csv.Values())
				}
			}
			if tc.Marshaled != "" {
				res, err := csv.MarshalJSON()
				if !assert.Nil(t, err) {
					t.Fatalf("Marshalling error: %v", err)
				}
				if !assert.Equal(t, tc.Marshaled, string(res)) {
					t.Fatalf("Marshalling error: %v", err)
				}
			}
		})
	}
}
