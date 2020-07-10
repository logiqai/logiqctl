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
Upload your LOGIQ deployment license
# logiqctl license set ./license.jws
# logiqctl license get 

`
var licenseLong = `
Logiq deployment comes configured with 30 day trial license
Obtain a valid license by contacting LOGIQ at license@logiq.ai.
This command helps set the subscription license for the deployment.
This command also help get the configured license information.
`

var licenseCmd = &cobra.Command{
	Use:     "license",
	Example: licenseExample,
	Aliases: []string{"licence"},
	Short:   "set and get license",
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
		Short:   "Retrive license information",
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
		Example: "logiqctl license set <license-file-path>",
		Aliases: []string{},
		Short:   "Configure license for LOGIQ deployment",
		PreRun:  utils.PreRun,
		Run: func(cmd *cobra.Command, args []string) {
			services.SetLicense(cmd, args)
		},
	}
	return cmd
}
