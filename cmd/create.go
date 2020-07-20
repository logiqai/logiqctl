package cmd

import (
	"github.com/logiqai/logiqctl/services"
	"github.com/logiqai/logiqctl/ui"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <resource_name>",
	Short: "Create a resource",
	Long: `Creates a reource from a resource specification. For example:

# Create a dashboard
logiqctl create dashboard -f <path to dashboard_spec_file.json>

# create eventrules
logiqctl create eventrules -f <path to eventrules_file.json>
`,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(ui.NewDashboardCreateCommand())
	createCmd.AddCommand(services.NewCreateEventRulesCommand())
}