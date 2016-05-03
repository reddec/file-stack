// Copyright © 2016 RedDec
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
	"fmt"
	"os"

	"github.com/reddec/file-stack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var stackFile string
var stack *fstack.Stack
var asJSON, asJSONbin bool
var msgSep string
var sep string

// This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "fstack",
	Short: "File base stack",
	Long:  `Command line utility to operate with file based stack`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func showMessage(headers, body []byte, notFirst bool) {
	if notFirst {
		msgSep = normalizeSeparator(msgSep)
		os.Stdout.Write([]byte(msgSep))
		os.Stderr.Write([]byte(msgSep))
	}
	var h map[string]string
	err := json.Unmarshal(headers, &h)
	if err != nil {
		panic(err)
	}
	if asJSONbin {
		msg := struct {
			Headers map[string]string `json:"headers"`
			Body    []byte            `json:"body"`
		}{}
		msg.Headers = h
		msg.Body = body
		data, _ := json.MarshalIndent(msg, "", "    ")
		os.Stdout.Write(data)
	} else if asJSON {
		msg := struct {
			Headers map[string]string `json:"headers"`
			Body    string            `json:"body"`
		}{}
		msg.Headers = h
		msg.Body = string(body)
		data, _ := json.MarshalIndent(msg, "", "    ")
		os.Stdout.Write(data)
	} else {
		sep = normalizeSeparator(sep)
		var second bool
		for k, v := range h {
			if second {
				fmt.Fprint(os.Stderr, sep)
			}
			second = true
			fmt.Fprintf(os.Stderr, "%s=%s", k, v)
		}
		os.Stdout.Write(body)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fstack.yaml)")
	RootCmd.PersistentFlags().StringVarP(&stackFile, "file", "f", "file.stack", "stack file name")

	RootCmd.PersistentFlags().BoolVarP(&asJSON, "json", "j", false, `output as json with string body`)
	RootCmd.PersistentFlags().BoolVar(&asJSONbin, "json-bin", false, `output as json with base64 body (replaces json)`)
	RootCmd.PersistentFlags().StringVarP(&sep, "delimiter", "d", "\\n", `delimiter between headers.
		Special names: \n, \0, \1`)
	RootCmd.PersistentFlags().StringVarP(&msgSep, "separator", "s", "\\n", `delimiter between messages (printed to stdout and stderr).
		Special names: \n, \0, \1`)

	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".fstack") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	fs, err := fstack.OpenStack(stackFile)
	if err != nil {
		panic(err)
	}
	stack = fs
}
