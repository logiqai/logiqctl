package cmd

import (
	"fmt"
	"os"

	"github.com/logiqai/logiqctl/ui"
	"github.com/logiqai/logiqctl/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <resource_name>",
	Short: "Delete one or more of your LOGIQ resources",
	Example: `
Delete dashboards
logiqctl delete dashboard slug		
	`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteDashboardsCommand())
}

func deleteDashboardsCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "dashboard",
		Example: "logiqctl delete dashboard slug",
		Aliases: []string{"dashboard"},
		Short:   "Delete LOGIQ Dashboards",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(os.Args) < 4 {
				fmt.Println("Slug not provided")
				os.Exit(1)
			}

			slug := os.Args[3]

			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Do you really want to delete dashboard: %s", slug),
				IsConfirm: true,
			}

			_, err := prompt.Run()
			if err != nil {
				fmt.Println("Exited without deleting dashboard")
				os.Exit(1)
			}

			if err := ui.DeleteDashboard(slug); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Successfully deleted dashboard: %s\n", slug)
		},
	}
}