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
	"github.com/logiqai/logiqctl/services"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/cobra"
)

var licenseExample = `
Upload your LOGIQ platform license
- logiqctl license set -f license.jws

View your LOGIQ license information
 - logiqctl license get 
`
var licenseLong = `
The LOGIQ Observability platform comes preconfigured with a 30-day trial license. You can obtain a valid license by contacting LOGIQ at license@logiq.ai.
This command lets you view your existing LOGIQ license or apply a new one. 
`

var licenseCmd = &cobra.Command{
	Use:     "license",
	Example: licenseExample,
	Aliases: []string{"licence"},
	Short:   "View or update LOGIQ license",
	Long:    licenseLong,
}

func init() {
	rootCmd.AddCommand(licenseCmd)
	licenseCmd.AddCommand(NewGetLicenseCommand())
	licenseCmd.AddCommand(NewSetLicenseCommand())
}

func NewGetLicenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "get",
		Example: "logiqctl license get",
		Aliases: []string{},
		Short:   "View license information",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.GetLicense()
		},
	}
	return cmd
}

func NewSetLicenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Example: "logiqctl license set -f <license-file-path>",
		Aliases: []string{},
		Short:   "Configure license for LOGIQ",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.SetLicense()
		},
	}
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	return cmd
}
