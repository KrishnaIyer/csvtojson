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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	conf "github.com/KrishnaIyer/csvtojson/pkg/config"
	"github.com/KrishnaIyer/csvtojson/pkg/csv"
	"github.com/KrishnaIyer/csvtojson/pkg/zephyrus"
	"github.com/spf13/cobra"
)

var (
	config = new(conf.Config)

	manager *conf.Manager

	// Root is the root of the commands.
	Root = &cobra.Command{
		Use:           "csvtojson",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "csvtojson is a simple command line tool to parse CSV files and convert them to JSON",
		Long:          `csvtojson is a simple command line tool to parse CSV files and convert them to JSON. More documentation at https://github.com/KrishnaIyer/csvtojson`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			err := manager.Unmarshal(config)
			if err != nil {
				panic(err)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logger, err := zephyrus.New(context.Background(), config.Debug)
			if err != nil {
				log.Fatal(err.Error())
			}
			defer logger.Clean()
			ctx := zephyrus.NewContextWithLogger(context.Background(), logger)

			raw, err := ioutil.ReadFile(config.CSVFile)
			if err != nil {
				logger.Fatal(err.Error())
			}

			csv, err := csv.New(ctx, raw, config.Parse)
			if err != nil {
				logger.Fatal(err.Error())
			}

			marshaled, err := json.Marshal(csv.Values())
			if err != nil {
				logger.Fatal(err.Error())
			}

			fmt.Println(string(marshaled))
		},
	}
)

// Execute ...
func Execute() {
	if err := Root.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}

func init() {
	manager = conf.New("config", *config)
	Root.PersistentFlags().AddFlagSet(manager.Flags())
}