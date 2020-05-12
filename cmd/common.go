/*
Copyright © 2020 Logiq.ai <cli@logiq.ai>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"

	"github.com/logiqai/logiqctl/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func preRun(cmd *cobra.Command, args []string) {

	cluster := viper.GetString(utils.KeyCluster)
	if cluster == "" {
		fmt.Println("Cluster is not set, run logiqctl config set-cluster END-POINT ")
		os.Exit(1)
	}
}

func preRunWithNs(cmd *cobra.Command, args []string) {

	cluster := viper.GetString(utils.KeyCluster)
	if cluster == "" {
		fmt.Println("Cluster is not set, run logiqctl config set-cluster END-POINT ")
		os.Exit(1)
	}

	ns := viper.GetString(utils.KeyNamespace)
	if ns == "" {
		fmt.Println("Context is not set, run logiqctl config set-context")
		os.Exit(1)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error Occured: %s", err.Error())
		os.Exit(-1)
	}
}
