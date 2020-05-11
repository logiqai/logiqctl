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
