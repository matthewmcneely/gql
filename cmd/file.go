// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/slothking-online/gql/client"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type FileCommandConfig struct {
	Config
	QueryFile     string
	VariablesFile string
}

func variablesFileFlag(variablesFile *string, flags *pflag.FlagSet) {
	flags.StringVar(variablesFile, "variables-file", "", "variables file")
}

func queryFileFlag(queryFile *string, flags *pflag.FlagSet) {
	flags.StringVar(queryFile, "query-file", "", "query file")
	cobra.MarkFlagRequired(flags, "query-file")
}

// rawCmd represents the raw command
func NewFileCommand(config FileCommandConfig) *cobra.Command {
	var (
		queryFile     string
		variablesFile string
		Endpoint      string
	)
	header := make(Header)
	fileCmd := &cobra.Command{
		Use:   "file",
		Short: "Execute graphql query from a file",
		Long:  `Executes query (with optional variables) located in a file or http endpoint against http GraphQL backend.`,
		Run: func(cmd *cobra.Command, args []string) {

			query, err := getFile(queryFile)
			if err != nil {
				log.Println("error getting contents of query-file", err)
				os.Exit(1)
			}
			if Endpoint == "" {
				if os.Getenv("ENDPOINT") != "" {
					Endpoint = os.Getenv("ENDPOINT")
				}
			}
			if Endpoint == "" {
				log.Println("GraphQL endpoint must be supplied via the --endpoint flag or environment variable ENDPOINT")
				os.Exit(1)
			}
			if variablesFile != "" {
				vars, err := getFile(variablesFile)
				if err != nil {
					log.Println("error getting contents of variables-file", err)
					os.Exit(1)
				}
				err = json.Unmarshal(vars, &variables)
				if err != nil {
					log.Println("error unmarshaling variables-file", err)
					os.Exit(1)
				}
			}
			httpHeader := make(http.Header)
			for k, v := range header {
				httpHeader.Add(k, v)
			}
			r := client.Raw{
				Query:         string(query),
				Variables:     map[string]interface{}(variables),
				OperationName: operationName,
				Header:        httpHeader,
			}
			cli := client.New(client.Config{
				Endpoint: Endpoint,
			})
			execute(config.Config, cli, r, nil)
		},
	}
	endpointFlag(&Endpoint, fileCmd.Flags())
	variablesFileFlag(&variablesFile, fileCmd.Flags())
	queryFileFlag(&queryFile, fileCmd.Flags())
	formatFlag(fileCmd.Flags())
	headersFlag(header, fileCmd.Flags())
	fileCmd.PersistentFlags().Var(
		variables,
		"set",
		"set grapqhl query variable, can be set multiple times",
	)
	fileCmd.PersistentFlags().StringVar(
		&operationName,
		"operation-name",
		"",
		"grapqhl operation name if provided query has more than one operation defined",
	)
	return fileCmd
}

func getFile(source string) ([]byte, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return nil, errors.Wrapf(err, "Error getting file contents from source %s", source)
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrapf(err, "Error reading file contents from source %s", source)
		}
		return data, nil
	}
	return ioutil.ReadFile(source)
}
