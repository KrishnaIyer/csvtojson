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
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/KrishnaIyer/csvtojson/pkg/zephyrus"
	"gopkg.in/yaml.v2"
)

var (
	errEmptyKey      = errors.New("Empty Key")
	errEmptyCSV      = errors.New("Empty CSV File")
	errSearchReplace = errors.New("Invalid search-replace pattern")

	replaceSeparator = ","
)

// CSV is the CSV file.
type CSV struct {
	values []map[string]string
}

// Config are the configuration options for the CSV decoder.
type Config struct {
	AllowMalformed bool   `name:"allow-malformed" description:"allow parsing malformed CSV"`
	FillEmptyWith  string `name:"fill-empty-with" description:"value to fill empty cells with. --allow-malformed must be set for this to be effective"`
	ReplaceWith    string `name:"replace-with" description:"simple text find and replace. Usage --replace-with <search>,<replacement>"`
}

// New parses the byte slice and creates a new CSV object.
func New(ctx context.Context, raw []byte, config Config) (*CSV, error) {
	logger := zephyrus.NewLoggerFromContext(ctx)
	var r *regexp.Regexp
	var replacement string
	replaceWith := config.ReplaceWith
	if replaceWith != "" {
		s := strings.Split(replaceWith, replaceSeparator)
		if len(s) != 2 {
			return nil, errSearchReplace
		}
		var err error
		pattern := s[0]
		logger.WithField("pattern", pattern).Info("Using search pattern")
		r, err = regexp.Compile(pattern)
		if err != nil {
			return nil, errSearchReplace
		}
		replacement = s[1]
	}
	reader := csv.NewReader(strings.NewReader(string(raw)))
	if config.AllowMalformed {
		reader.FieldsPerRecord = -1 // This allows variable number of columns per row.
	}
	readValues, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(readValues) < 2 {
		return nil, errEmptyCSV
	}

	values := make([]map[string]string, 0)
	for i := 1; i < len(readValues); i++ {
		keys := readValues[0]
		value := make(map[string]string)
		for j := 0; j < len(keys); j++ {
			noOfcolumns := len(readValues[i])
			if j >= noOfcolumns {
				value[keys[j]] = config.FillEmptyWith
				continue
			}
			// Find and replace if enabled
			if r != nil {
				val := r.ReplaceAllString(readValues[i][j], replacement)
				if val != "" {
					value[keys[j]] = val
					continue
				}
			}
			value[keys[j]] = readValues[i][j]
		}
		values = append(values, value)
	}
	return &CSV{
		values: values,
	}, nil
}

// Values returns the map of key:value pairs from the parsed CSV.
func (csv *CSV) Values() [](map[string]string) {
	return csv.values
}

// MarshalJSON marshals the read CSV values into JSON.
func (csv *CSV) MarshalJSON() ([]byte, error) {
	return json.Marshal(csv.Values())
}

// MarshalYAML marshals the read CSV values into YAML.
func (csv *CSV) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(csv.Values())
}
