// Copyright © 2016 RedDec <net.dev@mail.ru>
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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func normalizeSeparator(sep string) string {
	switch sep {
	case "\\n":
		sep = "\n"
	case "\\0":
		sep = string([]byte{0})
	case "\\1":
		sep = string([]byte{1})
	}
	return sep
}

// popCmd represents the pop command
var popCmd = &cobra.Command{
	Use:   "pop",
	Short: "POP opertation for stack",
	Long:  `Get last message from top of stack, print it and remove. Headers are printed to Stderr, body to Stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		if stack.Depth() == 0 {
			fmt.Fprintln(os.Stderr, "stack is empty")
			os.Exit(1)
		}
		headers, body, err := stack.Pop()
		if err != nil {
			panic(err)
		}
		showMessage(headers, body, false)
	},
}

func init() {
	RootCmd.AddCommand(popCmd)

}
