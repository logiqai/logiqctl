package utils

import (
	"fmt"
	"os"
	"github.com/logiqai/logiqctl/loglerpart"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PreRun(cmd *cobra.Command, args []string) {
	cluster := viper.GetString(KeyCluster)
	if cluster == "" {
		fmt.Println(`Cluster is not set, run "logiqctl config set-cluster END-POINT"" `)
		os.Exit(1)
	}
}

func PreRunUiTokenOrCredentials(cmd *cobra.Command, args []string) {
	PreRun(cmd, args)
	uiToken := viper.GetString(AuthToken)
	if uiToken == "" {
		user := viper.GetString(KeyUiUser)
		password := viper.GetString(KeyUiPassword)
		if user == "" && password == "" {
			fmt.Println(`Credentials must be set, run "logiqctl config help"`)
			os.Exit(1)
		}
	}
}

func PreRunWithNs(cmd *cobra.Command, args []string) {
	cluster := viper.GetString(KeyCluster)
	if cluster == "" {
		fmt.Println(`Cluster is not set, run "logiqctl config set-cluster END-POINT" `)
		os.Exit(1)
	}

	ns := viper.GetString(KeyNamespace)
	if ns == "" {
		fmt.Println("Context is not set, run logiqctl config set-context")
		os.Exit(1)
	} else {
		// fmt.Println("ns=", ns)
	}
}
func HandleError(err error) {
	if err != nil {
		if FlagEnablePsmod {
			loglerpart.DumpCurrentPsStat("ps_stat")
		}
		fmt.Printf("Err> %s\n", err.Error())
		os.Exit(1)
	}
}

func HandleError2(err error, mesg string) {
	if err != nil {
		if FlagEnablePsmod {
			loglerpart.DumpCurrentPsStat("ps_stat")
		}
		fmt.Printf("Err> %s\n     %s\n", mesg, err.Error())
		os.Exit(1)
	}
}
