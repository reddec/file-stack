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
	"io"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// tailCmd represents the tail command
var tailCmd = &cobra.Command{
	Use:   "tail",
	Short: "First N messages",
	Long:  `Get but not remove N first messages. Headers to Stderr, body to Stdout`,
	Run: func(cmd *cobra.Command, args []string) {
		var n uint64
		if count == 0 {
			return
		}
		err := stack.IterateForward(func(depth int, header io.Reader, body io.Reader) bool {
			hdata, err := ioutil.ReadAll(header)
			if err != nil {
				panic(err)
			}
			bdata, err := ioutil.ReadAll(body)
			if err != nil {
				panic(err)
			}
			showMessage(hdata, bdata, n > 0)
			n++
			return n < count
		})
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(tailCmd)
	tailCmd.PersistentFlags().Uint64VarP(&count, "count", "n", 10, "messages count")

}
