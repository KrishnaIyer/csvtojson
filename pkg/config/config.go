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

// Package config provides configuration functions.
package config

import (
	"fmt"
	"reflect"

	"github.com/KrishnaIyer/csvtojson/pkg/csv"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	commaSeparator = ","
	colonSeparator = ":"
)

// Config represents the configuration
type Config struct {
	CSVFile string `name:"csv-file" short:"c" description:"input csv file name"`
	OutFile string `name:"out-file" short:"o" description:"output file name. Use yaml or json based on the required format."`
	Debug   bool   `name:"debug" short:"d" description:"print detailed logs for errors"`
	YAML    bool   `name:"yaml" description:"marshal to yaml instead of json"`
	Parse   csv.Config
}

// Manager is the configuration manager.
type Manager struct {
	name  string
	flags *pflag.FlagSet
	viper *viper.Viper
}

// New returns a new initialized manager with the given config.
func New(name string, cfg Config) *Manager {
	flags := pflag.NewFlagSet(name, pflag.ExitOnError)
	viper := viper.New()
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	viper.SetConfigName(name)

	rootStruct := reflect.TypeOf(cfg)

	if rootStruct.Kind() != reflect.Struct {
		panic("configuration is not a struct")
	}

	// Parse the root struct separately
	flags.AddFlagSet(parseStructToFlags(rootStruct))

	// Find embedded structs and parse them. This currently only works one level deep.
	for i := 0; i < rootStruct.NumField(); i++ {
		fieldT := rootStruct.Field(i).Type
		if fieldT.Kind() == reflect.Struct {
			flags.AddFlagSet(parseStructToFlags(fieldT))
		}
	}

	err := viper.BindPFlags(flags)
	if err != nil {
		panic(err)
	}

	return &Manager{
		name:  name,
		flags: flags,
		viper: viper,
	}
}

// Flags returns pflag.FlagSet.
func (mgr *Manager) Flags() *pflag.FlagSet {
	return mgr.flags
}

// Unmarshal unmarshals the read config into the provided struct.
func (mgr *Manager) Unmarshal(config interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "name",
		Result:  config,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(mgr.viper.AllSettings())
}

// Viper returns viper.
func (mgr *Manager) Viper() *viper.Viper {
	return mgr.viper
}

// parseStructToFlags parses a struct and returns a flagset.
// It panics if there are parsing errors.
func parseStructToFlags(strT reflect.Type) *pflag.FlagSet {
	flags := pflag.FlagSet{}
	for i := 0; i < strT.NumField(); i++ {
		field := strT.Field(i)
		name := field.Tag.Get("name")
		if name == "" || name == "-" {
			continue
		}

		desc := field.Tag.Get("description")
		short := field.Tag.Get("short")

		fieldKind := field.Type.Kind()
		switch fieldKind {
		case reflect.String:
			flags.StringP(name, short, "", desc)
		case reflect.Bool:
			flags.BoolP(name, short, false, desc)
		default:
			panic(fmt.Errorf("Unknown type in config: %v", fieldKind))
		}
	}
	return &flags
}
