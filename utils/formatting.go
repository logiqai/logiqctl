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

package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"

	"gopkg.in/yaml.v2"

	"github.com/dustin/go-humanize"
)

var FlagOut string
var FlagTimeFormat string
var FlagNamespace string
var FlagCluster string
var FlagLogsSince string
var FlagLogsPageSize uint32
var FlagLogsFollow bool
var FlagProcId string
var FlagFile string
var FlagMaxFileSize int
var EventRuleGroupsFlag string

const LineBreaksKey = "lineBreaksAfterEachLogEntry"

func GetTimeAsString(s int64) string {
	t := time.Unix(s, 0)
	switch FlagTimeFormat {
	case "epoch":
		return fmt.Sprintf("%d", s)
	case "RFC3339":
		return t.Format(time.RFC3339)
	default:
		return humanize.Time(t)
	}
}

func PrintResponse(data interface{}) bool {
	if FlagOut == "json" {
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		return true
	} else if FlagOut == "yaml" {
		b, err := yaml.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
		return true
	}
	return false
}

func GetStartTime(lastSeen int64) time.Time {
	return time.Unix(lastSeen, 0).UTC()
}

func NeedsLineBreak() bool {
	if val, ok := viper.Get(LineBreaksKey).(bool); ok {
		return val
	}
	return false
}

func GetPageSize() uint32 {
	return FlagLogsPageSize
}
