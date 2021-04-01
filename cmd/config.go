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
	b64 "encoding/base64"
	"errors"
	"fmt"

	"github.com/manifoldco/promptui"

	"github.com/logiqai/logiqctl/services"

	"github.com/logiqai/logiqctl/utils"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config SUBCOMMAND",
	Short: "Modify your logiqctl configuration.",
	Long: `
The 'logiqctl config' command lets you configure your LOGIQ CLI. If this is your first time configuring logiqctl, you'll need an API token in order to use this command. To know how to generate an API token, read https://docs.logiq.ai/vewing-logs/logiqctl/obtaining-api-key.

Note: The values you provide during configuration will be written to the configuration file located at (~/.logiqctl)
`,
	Example: `
View current context
	logiqctl config view

Runs an interactive prompt that lets you configure logiqctl
	logiqctl config init

Set default cluster
	logiqctl config set-cluster END-POINT

Set default context
	logiqctl config set-context namespace

Runs an interactive prompt and lets you select a namespace from a list of namespaces
	logiqctl config set-context i

Set API token
	logiqctl config set-token api_token
`,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(NewSetClusterCommand())
	configCmd.AddCommand(NewSetContextCommand())
	configCmd.AddCommand(NewViewCommand())
	configCmd.AddCommand(NewSetConfigInitCommand())
	configCmd.AddCommand(NewCredentialsCommand())
	configCmd.AddCommand(NewUiTokenCommand())
}

func validInput(s string) error {
	if len(s) < 1 {
		return errors.New("cannot be empty")
	}
	return nil
}

func NewSetConfigInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init",
		Example: "logiqctl config init",
		Short:   "Configure logiqctl interactively",
		Run: func(cmd *cobra.Command, args []string) {
			clusterPrompt := promptui.Prompt{
				Label:    "Enter cluster endpoint ",
				Validate: validInput,
			}
			cluster, err := clusterPrompt.Run()
			if err != nil {
				fmt.Println("cannot read input")
				return
			}
			viper.Set(utils.KeyCluster, cluster)
			viper.Set(utils.KeyPort, utils.DefaultPort)

			err = services.Ping()
			if err != nil {
				return
			}

			userTokenPrompt := promptui.Prompt{
				Label:    "Enter User Token ",
				Validate: validInput,
				Mask:     '*',
			}

			userToken, err := userTokenPrompt.Run()
			if err != nil {
				fmt.Println("cannot read input")
				return
			}
			viper.Set(utils.AuthToken, userToken)

			selectedNs, err := services.RunSelectNamespacePrompt(false)
			if err != nil {
				fmt.Printf("Incorrect usage")
				return
			}
			viper.Set(utils.KeyNamespace, selectedNs)
			err = viper.WriteConfig()
			if err != nil {
				fmt.Print(err)
				return
			}
		},
	}
	return cmd
}

func NewCredentialsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set-credential",
		Example: "logiqctl set-credential login password",
		Short:   "Set your LOGIQ user credentials",
		Long: `
This command lets you set your LOGIQ user credentials. You'll need valid user credentials in order to access all operations.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				fmt.Println("Incorrect Usage")
				fmt.Println(cmd.Example)
				return
			}
			viper.Set(utils.KeyUiUser, b64.StdEncoding.EncodeToString([]byte(args[0])))
			viper.Set(utils.KeyUiPassword, b64.StdEncoding.EncodeToString([]byte(args[1])))
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
		Use:     "set-token",
		Example: "logiqctl set-token api_token",
		Short:   "Set your LOGIQ API token",
		Long: `
This command lets you set your LOGIQ API token. You'll need a valid LOGIQ cluster endpoint in order to complete all operations. 
		`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Println("Incorrect Usage")
				fmt.Println(cmd.Example)
				return
			}
			viper.Set(utils.AuthToken, args[0])
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
		Short:   "Set your LOGIQ platform endpoint",
		Long: `
This command lets you set your LOGIQ cluster endpoint. You'll need a valid LOGIQ cluster endpoint in order to complete all operations. 
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
		Short: "View your current logiqctl configuration.",
		Long: `
This command lets you view your current logiqctl configuration. 
		`,
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
		Example: "set-context <namespace name>",
		Short:   "Sets the default context or namespace.",
		Long: `
All logiqctl operations require a context or namespace. To override the default namespace set for an individual command, use the flag '-n'. 
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
		Short:   `Run an interactive prompt that lets you select a namespace from a list.`,
		Long: `
This command lets you set a default context interactively. Running 'logiqctl config set-context interactive' brings up an interactive list of namespaces from which you can select a namespace and set a context.
		`,
		Run: func(cmd *cobra.Command, args []string) {
			selectedNs, err := services.RunSelectNamespacePrompt(false)
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
	uiToken := viper.GetString(utils.AuthToken)
	if uiToken != "" {
		fmt.Printf("UI token set to: %s\n", uiToken)
	} else {
		fmt.Println("Default UI token is not set")
	}
}

func printUiCredentials() {
	uiUser := utils.GetUIUser()
	if uiUser != "" {
		fmt.Printf("UI user set to: %s\n", uiUser)
	} else {
		fmt.Println("Default UI user is not set")
	}
	uiPass := utils.GetUIPass()
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
