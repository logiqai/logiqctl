package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"log"
	"os"
	"strings"

	"github.com/logiqai/easymap"
	"github.com/logiqai/logiqbox/api/v1/realtimeLogStream"
	"github.com/logiqai/logiqbox/cfg"
	"google.golang.org/grpc"
)

var (
	matchMap       = map[string]interface{}{}
	allowAll       = false
	tailApps       = false
	tailNamespaces = false
	tailLabels     = false
	tailProcs      = false
)

func isMatch(logMap map[string]interface{}) bool {
	match := false
	if allowAll {
		return true
	}
	if tailApps {
		if _, found := matchMap[logMap["app_name"].(string)]; !found {
			return false
		} else {
			match = true
		}
	}

	if tailNamespaces {
		if ns, ok := logMap["namespace"]; ok {
			if _, matchOk := matchMap[ns.(string)]; !matchOk {
				return false
			} else {
				match = true
			}
		} else {
			if mRaw, ok := logMap["message_raw"]; ok {
				var v easymap.EasyMap
				if jErr := json.Unmarshal(([]byte)(mRaw.(string)), &v); jErr == nil {
					nsDetails := v.Get("namespace_name")
					if nsName, nsOk := nsDetails["kubernetes.namespace_name"]; nsOk {
						if _, matchOk := matchMap[nsName.(string)]; !matchOk {
							return false
						} else {
							match = true
						}
					}
				}
			} else {
				return false
			}
		}
	}

	if tailLabels {
		if mRaw, ok := logMap["message_raw"]; ok {
			var v easymap.EasyMap
			if jErr := json.Unmarshal(([]byte)(mRaw.(string)), &v); jErr == nil {
				labelsLookup := v.Get("labels")
				if labels, nsOk := labelsLookup["kubernetes.labels"]; nsOk {
					for k, v := range labels.(map[string]interface{}) {
						if matchValue, found := matchMap[k]; found {
							if matchValue.(string) == v.(string) {
								match = true
							}
						} else {
							return false
						}
					}
				} else {
					return false
				}
			}
		} else {
			return false
		}
	}
	return match
}

func Tail(c *cli.Context, config *cfg.Config, tA, tL, tN bool, apps []string) {
	tailApps = tA
	tailNamespaces = tN
	tailLabels = tL
	output := c.String("output")
	if !tA && !tL && !tN {
		allowAll = true
	} else {

		for _, match := range apps {
			if tailLabels {
				if len(strings.Split(match, ":")) != 2 {
					fmt.Println("label matches must use key:value format")
					os.Exit(-1)
				}
				lKey, lValue := strings.Split(match, ":")[0], strings.Split(match, ":")[1]
				matchMap[lKey] = lValue
			} else {
				matchMap[match] = true
			}
			if match == "*" {
				allowAll = true
			}
		}
	}
	conn, err := grpc.Dial(config.Cluster, grpc.WithInsecure())
	if err != nil {
		handleError(config, err)
		return
	}
	defer conn.Close()
	client := realtimeLogStream.NewLogStreamerServiceClient(conn)
	sub := &realtimeLogStream.Subscription{
		Applications: apps,
	}
	stream, err := client.StreamLog(context.Background(), sub)
	if err != nil {
		handleError(config, err)
		return
	}
	if output == OUTPUT_COLUMNS {
		printSyslogHeader()
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			handleError(config, err)
			return
		}
		logMap := make(map[string]interface{})

		err = json.Unmarshal([]byte(response.Log), &logMap)
		if err != nil {
			log.Print("Cannot read payload, this should not happen!")
		}

		ns := "default"
		if v, ok := logMap["msg_id"]; ok {
			ns = v.(string)
		}
		logMap["namespace"] = ns
		if isMatch(logMap) {
			printSyslogMessage(logMap, output)
		}

	}
}

//map[app_name:CRON facility:10 tag:CRON Id:bf4ef224-dc36-450b-a1d2-5909002f2aaf facility_string:security/authorization hostname:work severity:6 severity_string:info timestamp:2019-06-22T14:43:01.814112+05:30 uuid4:0ceacea4-600f-433d-ace4-7722401583fb message:pam_unix(cron:session): session opened for user tito by (uid=0) priority:86 proc_id:24280 source_ip:127.0.0.1]
