package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/manifoldco/promptui"

	"github.com/logiqai/logiqctl/utils"

	"github.com/logiqai/logiqctl/api/v1/query"
)

type templateType int

var templateTypeProcess templateType = 0
var templateTypeApplication templateType = 1
var once sync.Once

func GetTemplateForType(tType templateType) *promptui.SelectTemplates {
	var template = &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U000000BB {{ .Name | green }} ({{ .Details | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Details | red }})",
	}
	switch tType {
	case templateTypeProcess:
		template.Selected = "Process {{ .Name | red }} selected"
	case templateTypeApplication:
		template.Selected = "Application {{ .Name | red }} selected"
	default:
		template.Selected = "{{ .Name | red }} selected"
	}
	return template
}

type SelectDisplay struct {
	Name    string
	Details string
}

const (
	FMT = "%-24s | %-9s | %-16s | %-16s | %-16s | %s\n"
)

func printSyslogHeader() {
	fmt.Printf(FMT, "Timestamp", "Level", "Namespace", "Application", "Process/Pod", "Log message")
}

func printSyslogMessage(logMap map[string]interface{}, output string) {
	logMap["namespace"] = GetNamespaceSansHost(logMap["namespace"].(string))
	if output == OUTPUT_COLUMNS {
		fmt.Printf(FMT,
			logMap["timestamp"],
			logMap["severity_string"],
			logMap["namespace"],
			logMap["app_name"],
			logMap["proc_id"],
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
	} else {
		fmt.Printf("%s %9s %s %s %s %s",
			logMap["timestamp"],
			logMap["severity_string"],
			logMap["namespace"],
			logMap["app_name"],
			logMap["proc_id"],
			logMap["message"],
		)
		if !strings.HasSuffix(logMap["message"].(string), "\n") {
			fmt.Println()
		}
		if utils.NeedsLineBreak() {
			fmt.Println()
		}
	}
}

func PrintSyslogMessageForType(log *query.SysLogMessage, output string) {
	if output == OUTPUT_COLUMNS {
		fmt.Printf("%s | %s | %s | %s | %s\n",
			log.Timestamp,
			log.SeverityString,
			log.FacilityString,
			log.ProcID,
			log.Message,
		)
	} else if output == OUTPUT_RAW {
		fmt.Printf("%s %s %s %s - %s",
			log.Timestamp,
			log.SeverityString,
			log.FacilityString,
			log.ProcID,
			log.Message,
		)
		if !strings.HasSuffix(log.Message, "\n") {
			fmt.Println()
		}
		if utils.NeedsLineBreak() {
			fmt.Println()
		}
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

var hostnameSuffix = "-@h"

func GetNamespaceSansHost(namespace string) string {
	ns := namespace
	if namespace != "" && strings.HasSuffix(namespace, hostnameSuffix) {
		ns = strings.Replace(namespace, hostnameSuffix, "", -1)
	}
	return ns
}
