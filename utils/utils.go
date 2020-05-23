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
	"fmt"

	"github.com/spf13/viper"
)

const (
	KeyCluster    = "cluster"
	KeyPort       = "port"
	DefaultPort   = "8081"
	KeyNamespace  = "namespace"
	KeyUiToken    = "uitoken"
	KeyUiUser     = "uiuser"
	KeyUiPassword = "uipassword"
)

func GetClusterUrl() string {
	var cluster string
	if FlagCluster != "" {
		cluster = FlagCluster
	} else {
		cluster = viper.GetString(KeyCluster)
	}
	port := viper.GetString(KeyPort)
	return fmt.Sprintf("%s:%s", cluster, port)
}

func GetDefaultNamespace() string {
	if FlagNamespace != "" {
		return FlagNamespace
	}
	ns := viper.GetString(KeyNamespace)
	return ns
}
