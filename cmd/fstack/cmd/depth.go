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

	"github.com/spf13/cobra"
)

// depthCmd represents the depth command
var depthCmd = &cobra.Command{
	Use:   "depth",
	Short: "Stack depth",
	Long:  `Get messages count in stack`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(stack.Depth())
	},
}

func init() {
	RootCmd.AddCommand(depthCmd)
}
