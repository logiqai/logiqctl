package services

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/logiqai/logiqctl/api/v1/query"
	"github.com/logiqai/logiqctl/cfg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(config *cfg.Config, err error) {
	//fmt.Printf("%v\n",err)
	errorStatus := status.Convert(err)
	switch errorStatus.Code() {
	case codes.Unavailable:
		fmt.Printf("Could not connect, Are you sure you can connect to [%s]?", config.Cluster)
		fmt.Println("Please verify your configuration...")
		fmt.Println("++++++++++++++++++++++")
		fmt.Print(config.String())
		fmt.Println("++++++++++++++++++++++")
		os.Exit(1)
	case codes.PermissionDenied:
		fmt.Printf("You are not allowed to do this operations, this may require Administrator privileges\n")
		os.Exit(1)
	case codes.Unauthenticated:
		fmt.Printf("Invalid Token, You may need to Login\n")
		os.Exit(1)
	case codes.AlreadyExists:
		fmt.Printf("Error : %s\n", errorStatus.Message())
		os.Exit(1)
	case codes.NotFound:
		fmt.Printf("Error : %s\n", errorStatus.Message())
		os.Exit(1)
	case codes.InvalidArgument:
		fmt.Printf("Error : %s\n", errorStatus.Message())
		os.Exit(1)
	default:
		fmt.Printf("Error : [%s]\n", errorStatus.Message())
		break
	}
}

const (
	FMT = "%-24s | %-16s | %-16s | %-16s | %-16s | %s\n"
)

func printSyslogHeader() {
	fmt.Printf(FMT, "Timestamp", "Namespace", "Application", "Process/Pod", "Facility", "Log message")
}

func printSyslogMessage(logMap map[string]interface{}, output string) {
	if output == OUTPUT_COLUMNS {
		fmt.Printf(FMT,
			logMap["timestamp"],
			logMap["namespace"],
			logMap["app_name"],
			logMap["proc_id"],
			logMap["facility_string"],
			logMap["message"],
		)
	} else if output == OUTPUT_RAW {
		fmt.Printf("%s %s %s %s %s\n",
			logMap["timestamp"],
			logMap["namespace"],
			logMap["app_name"],
			logMap["proc_id"],
			logMap["facility_string"],
			logMap["message"],
		)
	} else if output == OUTPUT_JSON {
		v, err := json.Marshal(logMap)
		if err == nil {
			fmt.Printf("%s\n", v)
		} else {
			fmt.Printf("Error marshalling JSON %v", logMap)
			os.Exit(-1)
		}
	}
}

func printSyslogMessageForType(log *query.SysLogMessage, output string) {
	if output == OUTPUT_COLUMNS {
		fmt.Printf("%-28s|%-6s|%s|%s|%-5s|%s|%s\n",
			log.Timestamp,
			log.SeverityString,
			log.Hostname,
			log.ProcID,
			log.AppName,
			log.FacilityString,
			log.Message,
		)
	} else if output == OUTPUT_RAW {
		fmt.Printf("%s %s %s %s %s %s %s\n",
			log.Timestamp,
			log.SeverityString,
			log.FacilityString,
			log.Hostname,
			log.AppName,
			log.ProcID,
			log.Message,
		)
	} else if output == OUTPUT_JSON {
		v, err := json.Marshal(log)
		if err == nil {
			fmt.Printf("%s\n", v)
		} else {
			fmt.Printf("Error marshalling JSON %v", *log)
			os.Exit(-1)
		}
	}
}
