package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	KeyCluster   = "cluster"
	KeyPort      = "port"
	DefaultPort  = "8081"
	KeyNamespace = "namespace"
)

func GetClusterUrl() string {
	cluster := viper.GetString(KeyCluster)
	port := viper.GetString(KeyPort)
	return fmt.Sprintf("%s:%s", cluster, port)
}

func GetDefaultNamespace() string {
	ns := viper.GetString(KeyNamespace)
	return ns
}
