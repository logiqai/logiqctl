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

package services

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/logiqai/logiqctl/utils"

	"github.com/logiqai/easymap"
	"github.com/logiqai/logiqctl/api/v1/realtimeLogStream"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	matchAppMap       = map[string]bool{}
	matchNamespaceMap = map[string]bool{}
	matchLabelsMap    = map[string]string{}
	matchProcessMap   = map[string]bool{}

	allowAll       = false
	tailApps       = false
	tailNamespaces = false
	tailLabels     = false
	tailProcs      = false
	tailDefault    = false
)

func isMatch(logMap map[string]interface{}) bool {
	match := false

	log.Debugln("isMatch:", tailApps, tailNamespaces, tailProcs, tailLabels)
	log.Debugln("isMatch:", logMap["app_name"], logMap["proc_id"])
	log.Debugln("")

	if allowAll {
		return true
	}

	if tailApps {
		if mApp, ok := logMap["app_name"]; ok {
			if _, found := matchAppMap[mApp.(string)]; !found {
				return false
			} else {
				log.Debugln("Matched app_name ", mApp)
				match = true
			}
		} else {
			log.Debugln("No app_name is data", logMap)
		}
	}

	if tailNamespaces {
		if ns, ok := logMap["namespace"]; ok {
			if _, matchOk := matchNamespaceMap[ns.(string)]; !matchOk {
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
						if _, matchOk := matchNamespaceMap[nsName.(string)]; !matchOk {
							return false
						} else {
							log.Debugln("Matched namespace ", ns)
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
						if matchValue, found := matchLabelsMap[k]; found {
							if matchValue == v.(string) {
								log.Debugln("Matched label k:", k, " v:", v)
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

	if tailProcs {
		if mProc, ok := logMap["proc_id"]; ok {
			if _, found := matchProcessMap[mProc.(string)]; !found {
				return false
			} else {
				log.Debugln("Matched process ", mProc)
				match = true
			}
		} else {
			log.Debugln("No proc_id in data", logMap)
		}
	}

	return match
}

func setupMatchAttributeMaps(matches []string, m map[string]bool) {
	for _, v := range matches {
		m[v] = true
	}
}

func setupMatchAttributeValueMaps(matches []string, sep string, m map[string]string) {
	log.Debugln(matches)
	for _, v := range matches {
		sp := strings.Split(v, sep)
		if len(sp) != 2 {
			log.Fatal("Labels matches must include key and value with a : or = separator")
		}
		m[sp[0]] = sp[1]
	}
}

func Tail(tN, tL, tA, tP, def []string) error {
	log.Debugln(len(tA), len(tN), len(tL), len(tP), len(def))
	tailApps = len(tA) > 0
	tailNamespaces = len(tN) > 0
	tailLabels = len(tL) > 0
	tailProcs = len(tP) > 0
	tailDefault = len(def) > 0
	output := utils.FlagOut
	log.Debugln(tN, tL, tA, tP, def)
	log.Debugln("A:", tailApps, "N:", tailNamespaces, "L:", tailLabels, "P:", tailProcs, "D:", tailDefault)
	if !tailApps && !tailLabels && !tailNamespaces && !tailProcs && !tailDefault {
		allowAll = true
	}
	setupMatchAttributeMaps(tA, matchAppMap)
	setupMatchAttributeMaps(tN, matchNamespaceMap)
	setupMatchAttributeMaps(tP, matchProcessMap)
	setupMatchAttributeValueMaps(tL, ":=", matchLabelsMap)

	log.Debugln(matchNamespaceMap, matchLabelsMap, matchProcessMap, matchAppMap)

	conn, err := grpc.Dial(utils.GetClusterUrl(), grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := realtimeLogStream.NewLogStreamerServiceClient(conn)
	if tA == nil || len(tA) == 0 {
		//always default to all logs
		tA = []string{"*"}
	}
	sub := &realtimeLogStream.Subscription{
		Applications: tA,
	}
	stream, err := client.StreamLog(context.Background(), sub)
	if err != nil {
		return err
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
			return err
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
	return nil
}

//map[app_name:CRON facility:10 tag:CRON Id:bf4ef224-dc36-450b-a1d2-5909002f2aaf facility_string:security/authorization hostname:work severity:6 severity_string:info timestamp:2019-06-22T14:43:01.814112+05:30 uuid4:0ceacea4-600f-433d-ace4-7722401583fb message:pam_unix(cron:session): session opened for user tito by (uid=0) priority:86 proc_id:24280 source_ip:127.0.0.1]
