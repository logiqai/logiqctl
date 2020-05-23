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

	"github.com/logiqai/logiqctl/services"

	"github.com/logiqai/logiqctl/utils"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config SUBCOMMAND",
	Short: "Modify logiqctl configuration options",
	Long: `
# View current context
	logiqctl config view

# Set default cluster
	logiqctl set-cluster END-POINT

# Set default context
	logiqctl set-context namespace

# Runs an interactive prompt and let user select namespace from the list
	logiqctl set-context i

# Set ui token context
	logiqctl set-ui-token token

# Set ui credential
	logiqctl set-ui-credential user password
`,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(NewSetClusterCommand())
	configCmd.AddCommand(NewSetContextCommand())
	configCmd.AddCommand(NewViewCommand())
	configCmd.AddCommand(NewUiTokenCommand())
	configCmd.AddCommand(NewUiCredentialsCommand())
}

func NewUiCredentialsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-ui-credential",
		Example: "logiqctl set-ui-credential login password",
		Short:   "Sets a logiq ui credentials",
		Long: `
Sets the cluster ui credentials, a valid logiq cluster end point is also required for all the operations
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Println("Incorrect Usage")
				fmt.Println(cmd.Example)
				return
			}
			viper.Set(utils.KeyUiUser, args[0])
			viper.Set(utils.KeyUiPassword, args[1])
			err := viper.WriteConfig()
			if err != nil {
				fmt.Print(err)
				return
			}
			printUiCredentials()
		},
	}
	return cmd
}

func NewUiTokenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-ui-token",
		Example: "logiqctl set-ui-token api_access_token",
		Short:   "Sets a logiq ui api token",
		Long: `
Sets the cluster UI api token, a valid logiq cluster end point is also required for all the operations
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Incorrect Usage")
				fmt.Println(cmd.Example)
				return
			}
			viper.Set(utils.KeyUiToken, args[0])
			err := viper.WriteConfig()
			if err != nil {
				fmt.Print(err)
				return
			}
			printUiToken()
		},
	}
	return cmd
}

func NewSetClusterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-cluster",
		Example: "logiqctl set-cluster END-POINT",
		Short:   "Sets a logiq cluster entry point",
		Long: `
Sets the cluster, a valid logiq cluster end point is required for all the operations
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Incorrect Usage")
				fmt.Println(cmd.Example)
				return
			}
			viper.Set(utils.KeyCluster, args[0])
			viper.Set(utils.KeyPort, utils.DefaultPort)
			err := viper.WriteConfig()
			if err != nil {
				fmt.Print(err)
				return
			}
			printCluster()
		},
	}
	return cmd
}

func NewViewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "View current defaults",
		Run: func(cmd *cobra.Command, args []string) {
			printCluster()
			printNamespace()
			printUiToken()
		},
	}
}

func NewSetContextCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-context",
		Example: "set-context NAMESPACE",
		Short:   "Sets a default namespace in logiqctl",
		Long: `
This will the default context for all the operations.
		`,
		PreRun: utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("Incorrect usage")
				fmt.Println(cmd.Example)
				return
			}
			if len(args) == 1 {
				setContext(args[0])
				printNamespace()
				return
			}
		},
	}
	cmd.AddCommand(NewInteractiveSetContextCommand())
	return cmd
}

func NewInteractiveSetContextCommand() *cobra.Command {
	var interactiveCmd = &cobra.Command{
		Use:     "interactive",
		Aliases: []string{"i"},
		Short:   `Runs an interactive prompt and let user select namespace from the list`,
		Run: func(cmd *cobra.Command, args []string) {
			selectedNs, err := services.RunSelectNamespacePrompt()
			if err != nil {
				fmt.Printf("Incorrect usage")
				fmt.Println(cmd.Example)
				return
			}
			setContext(selectedNs)
			printNamespace()
		},
	}
	return interactiveCmd
}

func setContext(arg string) {
	viper.Set(utils.KeyNamespace, arg)
	err := viper.WriteConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
}

func printUiToken() {
	uiToken := viper.GetString(utils.KeyUiToken)
	if uiToken != "" {
		fmt.Printf("UI token set to: %s\n", uiToken)
	} else {
		fmt.Println("Default UI token is not set")
	}
}

func printUiCredentials() {
	uiUser := viper.GetString(utils.KeyUiUser)
	if uiUser != "" {
		fmt.Printf("UI user set to: %s\n", uiUser)
	} else {
		fmt.Println("Default UI user is not set")
	}
	uiPass := viper.GetString(utils.KeyUiPassword)
	if uiPass != "" {
		fmt.Printf("UI password set to: %s\n", uiPass)
	} else {
		fmt.Println("Default UI password is not set")
	}
}

func printCluster() {
	cluster := viper.GetString(utils.KeyCluster)
	if cluster != "" {
		fmt.Printf("Cluster Endpoint set to: %s\n", cluster)
	} else {
		fmt.Println("Default Cluster not set")
	}
}

func printNamespace() {
	ns := viper.GetString(utils.KeyNamespace)
	if ns != "" {
		fmt.Printf("Default Context set to: %s\n", ns)
	} else {
		fmt.Println("Default Context not set")
	}
}
