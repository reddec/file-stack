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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var headers []string

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "PUSH opertation for stack",
	Long:  `Read all data from STDIN as single message and push to stack`,
	Run: func(cmd *cobra.Command, args []string) {
		heads := map[string]string{}
		for _, head := range headers {
			parts := strings.SplitN(head, "=", 2)
			if len(parts) != 2 {
				panic("BAD header: must be key=value")
			}
			heads[parts[0]] = parts[1]
		}
		body, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		data, _ := json.Marshal(heads)
		id, err := stack.Push(data, body)
		if err != nil {
			panic(err)
		}
		fmt.Println(id)
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
	pushCmd.PersistentFlags().StringSliceVarP(&headers, "header", "H", []string{}, "set headers")
}
