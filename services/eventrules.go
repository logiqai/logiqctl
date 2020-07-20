package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/logiqai/logiqctl/api/v1/eventRules"
	"github.com/logiqai/logiqctl/grpc_utils"
	"github.com/logiqai/logiqctl/utils"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// logiqctl create eventrules -f=filename.json
// logiqctl get eventrules all
// logiqctl get eventrules all -w=filename.json
// logiqctl get eventrules groups
// logiqctl get eventrules groups -g=group1,group2,...
// logiqctl get eventrules groups -g=group1,group2,... -w=filename.json

var eventRuleGroupsFlag string

func NewCreateEventRulesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eventrules",
		Example: "logiqctl create eventrules -f <path to event rules file>",
		Aliases: []string{"eventrule", "er"},
		Short:   "Import event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			if utils.FlagFile == "" {
				fmt.Println("Missing event rules file")
				return
			} else {
				fmt.Println("Event rules file :", utils.FlagFile)
				if fileBytes, err := ioutil.ReadFile(utils.FlagFile); err != nil {
					fmt.Println("Unable to read file", utils.FlagFile)
					return
				} else {
					rules := []eventRules.EventRule{}
					if err = json.Unmarshal(fileBytes, &rules); err != nil {
						fmt.Println("Unable to decode event rules from ", utils.FlagFile)
					} else {
						createEventRules(rules)
					}
				}
			}
		},
	}
	cmd.PersistentFlags().StringVarP(&utils.FlagFile, "file", "f", "", "Path to file")
	return cmd
}

func createEventRules(ers []eventRules.EventRule) error {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer conn.Close()
	client := eventRules.NewEventRulesServiceClient(conn)

	for _, er := range ers {
		if _, err := client.CreateEventRule(grpc_utils.GetGrpcContext(), &er); err != nil {
			fmt.Printf("Fail %s not created. Error: %s\n", er.Name, err.Error())
		} else {
			fmt.Printf("Success %s created\n", er.Name)
		}
	}
	return nil
}

func NewGetEventRulesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "eventrules",
		Example: "logiqctl get eventrule",
		Aliases: []string{"eventrule", "er"},
		Short:   "List event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			help := `Usage:
logiqctl get eventrules all
logiqctl get eventrules all -w filename.json
logiqctl get eventrules groups
logiqctl get eventrules groups -g=group1,group2,...
logiqctl get eventrules groups -g=group1,group2,... -w filename.json
`
			fmt.Print(help)
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "all",
		Example: "logiqctl get eventrules all",
		Aliases: []string{},
		Short:   "List all event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			getEventRules(args, nil)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "groups",
		Example: "logiqctl get eventrules all",
		Aliases: []string{},
		Short:   "List all event rules",
		PreRun:  utils.PreRunWithNs,
		Run: func(cmd *cobra.Command, args []string) {
			getEventRuleGroups(args)
		},
	})
	cmd.PersistentFlags().StringVarP(&eventRuleGroupsFlag, "groups", "g", "", "list of groups separated by comma")
	cmd.PersistentFlags().StringVarP(&utils.FlagFile, "file", "w", "", "Path to file to be written to")
	return cmd
}
func getEventRuleGroups(args []string) error {
	if eventRuleGroupsFlag != "" {
		// fmt.Println(eventRuleGroupsFlag)
		groups := strings.Split(eventRuleGroupsFlag, ",")
		if len(groups) == 0 {
			err := fmt.Errorf("Use logiqctl get eventrules groups -g=group1,group2,...")
			fmt.Println(err.Error())
			return err
		} else {
			fmt.Println("Fetching event rules for groups: ", groups)
			return getEventRules(args, groups)
		}
	} else {
		conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		defer conn.Close()
		client := eventRules.NewEventRulesServiceClient(conn)

		if groups, err := client.GetEventRuleGroups(grpc_utils.GetGrpcContext(), &emptypb.Empty{}); err != nil {
			fmt.Println(err.Error())
			return err
		} else {
			if len(groups.GroupNames) == 0 {
				fmt.Println("No event rules groups found")
				return nil
			} else {
				fmt.Println("Number of event rules groups found:", len(groups.GroupNames))
			}
			for _, g := range groups.GroupNames {
				fmt.Println(g)
			}
			return nil
		}
	}
}

func getEventRules(args, groupNames []string) error {
	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer conn.Close()

	in := &eventRules.EventRulesQueryRequest{Count: 1024}
	inFilter := &eventRules.EventRulesFilter{}
	if len(groupNames) > 0 {
		inFilter.GroupNames = groupNames
	}
	in.Filter = inFilter

	client := eventRules.NewEventRulesServiceClient(conn)
	if ers, err := client.GetEventRules(grpc_utils.GetGrpcContext(), in); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		if len(ers.EventRules) == 0 {
			fmt.Println("No event rules found")
		} else {
			fmt.Println("Number of event rules found:", len(ers.EventRules))
		}

		if eventRules, err := json.MarshalIndent(ers.EventRules, " ", "    "); err != nil {
			fmt.Println(err.Error())
			return err
		} else {
			if utils.FlagFile == "" {
				fmt.Println(string(eventRules))
				return nil
			} else {
				if err := ioutil.WriteFile(utils.FlagFile, eventRules, 0644); err != nil {
					fmt.Println(err.Error())
					return err
				} else {
					fmt.Printf("Event rules exported to %s \n", utils.FlagFile)
					return nil
				}
			}
		}
	}
}
