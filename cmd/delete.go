package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/logiqai/logiqctl/types"
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
	deleteCmd.AddCommand(deleteQueriesCommand())
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

			res, err := ui.GetDashboard([]string{slug})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			var dashboard types.DashboardSpec
			if err = json.Unmarshal([]byte(res), &dashboard); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			prompt.Label = fmt.Sprintf("Do you want to delete all the queries related to dashboard : %s", slug)
			_, err = prompt.Run()
			if err != nil {
				fmt.Println("Deleted dashboard but kept queries as-is")
				os.Exit(1)
			}

			for _, w := range dashboard.Widgets {
				if err := ui.DeleteQuery(strconv.Itoa(w.Visualization.Query.Id)); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}

			fmt.Printf("Successfully deleted dashboard: %s and it's related queries\n", slug)
		},
	}
}

func deleteQueriesCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "query",
		Example: "logiqctl delete query id",
		Aliases: []string{"query"},
		Short:   "Delete LOGIQ Queries",
		PreRun:  utils.PreRunUiTokenOrCredentials,
		Run: func(cmd *cobra.Command, args []string) {
			if len(os.Args) < 4 {
				fmt.Println("ID was not provided")
			}

			ID := os.Args[3]
			prompt := promptui.Prompt{
				Label:     fmt.Sprintf("Do you really want to delete query: %s", ID),
				IsConfirm: true,
			}

			_, err := prompt.Run()
			if err != nil {
				fmt.Println("Exited without deleting the query")
				os.Exit(1)
			}

			if err := ui.DeleteQuery(ID); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("Successfully delete the query: %s", ID)
		},
	}
}