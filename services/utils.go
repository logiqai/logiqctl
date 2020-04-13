package services

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/logiqai/logiqbox/api/v1/query"
	"github.com/logiqai/logiqbox/cfg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(config *cfg.Config, err error) {
	errorStatus := status.Convert(err)
	switch errorStatus.Code() {
	case codes.Unavailable:
		log.Printf("Could not connect, Are you sure you can connect to [%s]?", config.Cluster)
		log.Println("Please verify your configuration...")
		fmt.Println("++++++++++++++++++++++")
		fmt.Print(config.String())
		fmt.Println("++++++++++++++++++++++")
		break
	case codes.Unauthenticated:
		log.Println("Invalid Token, You may need to Login")
		break
	default:
		log.Printf("Error Code: [%d, %s]", errorStatus.Code(), errorStatus.Code().String())
		break
	}
}

const (
	FMT = "%-24s | %-16s | %-16s | %-16s | %s\n"
)
func printSyslogHeader() {
	fmt.Printf(FMT, "Timestamp","Application","Process/Pod","Facility","Log message")
}

func printSyslogMessage(logMap map[string]interface{}, output string) {
	if output == OUTPUT_COLUMNS {
		fmt.Printf(FMT,
			logMap["timestamp"],
			logMap["app_name"],
			logMap["proc_id"],
			logMap["facility_string"],
			logMap["message"],
		)
	} else if output == OUTPUT_RAW {
		fmt.Printf("%s %s %s %s %s\n",
			logMap["timestamp"],
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
