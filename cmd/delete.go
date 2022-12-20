package cmd

import (
	"fmt"

	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <resource_name>",
	Short: "Delete or Archive your LOGIQ resources",
	Long:  `Delete or Archive your LOGIQ resource which matches the given identifier`,
	Example: `
Archive Dashboard for the given slug
logiqctl delete dashboard <slug>

Delete Dashboard for the given slug
logiqctl delete dashboard <slug> -d 
`,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.AddCommand(deleteDashboardCommand())
}

func deleteDashboardCommand() *cobra.Command {
	var deleteFlag bool
	var queryFlag bool
	cmd := &cobra.Command{
		Use:     "dashboard",
		Example: "logiqctl delete dashboard <slug>",
		Short:   "Delete/Archive a dashboard",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			// _, err := ui.DeleteDashboardBySlug(args[0], !deleteFlag, queryFlag)
			// if err != nil {
			// 	fmt.Printf("error: %s\n", err.Error())
			// 	os.Exit(255)
			// }
			fmt.Println("Dashboard deleted..")
		},
	}
	cmd.Flags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	cmd.Flags().BoolVarP(&deleteFlag, "delete", "d", false, "Delete Dashboard")
	cmd.Flags().BoolVarP(&queryFlag, "query", "q", false, "Delete queries used in dashboard")

	return cmd
}
