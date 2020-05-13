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
	"path"

	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Short:   "Logiqctl - CLI for Logiq Observability stack",
	Use:     "logiqctl [flags] [options]",
	Version: "1.0.2",
	Long: `
The LOGIQ command line toolkit, logiqctl, allows you to run commands against LOGIQ Observability stack. 
- You can tail logs from your applications and servers
- View available namespaces/applications/processes
- Query historical data 
- Search your log data.
- View Events


Find more information at: https://docs.logiq.ai/logiqctl/logiq-box

`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&utils.FlagOut, "output", "o", "table", `Output format. One of: table|json|yaml.`)
	rootCmd.PersistentFlags().StringVarP(&utils.FlagTimeFormat, "time-format", "t", "relative", `Time formatting options. One of: relative|epoch|RFC3339. 
This is only applicable when the output format is table. json and yaml outputs will have time in epoch seconds.
json output is not indented, use '| jq' for advanced json operations`)
	rootCmd.PersistentFlags().StringVarP(&utils.FlagNamespace, "namespace", "n", "", "Override the default context set by `logiqctl set-context' command")
	rootCmd.PersistentFlags().StringVarP(&utils.FlagCluster, "cluster", "c", "", "Override the default cluster set by `logiqctl set-cluster' command")
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configDir := path.Join(home, ".logiqctl")
	exists, err := exists(configDir)
	cfgFile := path.Join(configDir, "logiqctl.toml")
	if err != nil {
		fmt.Println(fmt.Errorf("Cannot create config: %s ", err.Error()))
		return
	}
	if !exists {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			fmt.Println(fmt.Errorf("Cannot create config: %s ", err.Error()))
			return
		}
		viper.SetConfigFile(cfgFile)
		viper.Set("logiqctl", "v.1.0")
		viper.Set(utils.LineBreaksKey, false)
		viper.WriteConfig()
	} else {
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Cannot Use Config file:", viper.ConfigFileUsed())
	}
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
