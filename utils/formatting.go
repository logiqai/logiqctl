package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/dustin/go-humanize"
)

var FlagOut string
var FlagTimeFormat string

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
