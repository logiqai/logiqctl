package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/logiqai/logiqctl/api/v1/eventRules"
	"github.com/logiqai/logiqctl/grpc_utils"
	"github.com/logiqai/logiqctl/utils"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// logiqctl create eventrules -f=filename.json
// logiqctl get eventrules all
// logiqctl get eventrules all -w=filename.json
// logiqctl get eventrules groups
// logiqctl get eventrules groups -g=group1,group2,...
// logiqctl get eventrules groups -g=group1,group2,... -w=filename.json

func CreateEventRules(ers []eventRules.EventRule) error {
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

func GetEventRuleGroups(args []string) error {
	if utils.EventRuleGroupsFlag != "" {
		// fmt.Println(eventRuleGroupsFlag)
		groups := strings.Split(utils.EventRuleGroupsFlag, ",")
		if len(groups) == 0 {
			err := fmt.Errorf("Use logiqctl get eventrules groups -g=group1,group2,...")
			fmt.Println(err.Error())
			return err
		} else {
			fmt.Println("Fetching event rules for groups: ", groups)
			return GetEventRules(args, groups)
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

func GetEventRules(args, groupNames []string) error {
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
