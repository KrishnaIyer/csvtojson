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

// Package csv parses CSV files.
package csv

import (
	"context"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"strings"
)

var errEmptyKey = errors.New("")

// CSV is the CSV file.
type CSV struct {
	values map[string][]string
}

// New parses the file and creates a new CSV object.
func New(ctx context.Context, fileName string) (*CSV, error) {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(strings.NewReader(string(raw)))
	values, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	csv := &CSV{
		values: make(map[string][]string),
	}
	// Extract the values into a map
	if len(values) > 0 {
		keys := values[0]
		for i, row := range values {
			if i == 0 {
				continue
			}
			for j, val := range row {
				key := keys[j]
				if key == "" {
					return nil, error
				}
				csv.values[key] = append(csv.values[key], val)
			}
		}
	}
	return csv, nil
}
