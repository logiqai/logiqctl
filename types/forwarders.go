package types

import (
	"fmt"
)

type ForwarderConfig struct {
	ApiKey                         string `json:"api_key,omitempty" yaml:"api_key,omitempty" yaml:"api_key,omitempty" yaml:"api_key,omitempty"`
	Host                           string `json:"host,omitempty" yaml:"host,omitempty" yaml:"host,omitempty" yaml:"host,omitempty"`
	Tags                           string `json:"tags,omitempty" yaml:"tags,omitempty" yaml:"tags,omitempty" yaml:"tags,omitempty"`
	Type                           string `json:"type,omitempty" yaml:"type,omitempty" yaml:"type,omitempty" yaml:"type,omitempty"`
	BufferSize                     string `json:"buffer_size,omitempty" yaml:"buffer_size,omitempty" yaml:"buffer_size,omitempty" yaml:"buffer_size,omitempty"`
	Port                           string `json:"port,omitempty" yaml:"port,omitempty"`
	AccessKey                      string `json:"access_key,omitempty" yaml:"access_key,omitempty"`
	Bucket                         string `json:"bucket,omitempty" yaml:"bucket,omitempty"`
	Endpoint                       string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	FilterForward                  string `json:"filter_forward,omitempty" yaml:"filter_forward,omitempty"`
	FlushToBucketIntervalInSeconds string `json:"flush_to_bucket_interval_in_seconds,omitempty" yaml:"flush_to_bucket_interval_in_seconds,omitempty"`
	MaxObjectSizeInBytes           string `json:"max_object_size_in_bytes,omitempty" yaml:"max_object_size_in_bytes,omitempty"`
	Region                         string `json:"region,omitempty" yaml:"region,omitempty"`
	SecretKey                      string `json:"secret_key,omitempty" yaml:"secret_key,omitempty"`
	Stream                         string `json:"stream,omitempty" yaml:"stream,omitempty"`
	Password                       string `json:"password,omitempty" yaml:"password,omitempty"`
	User                           string `json:"user,omitempty" yaml:"user,omitempty"`
	AccountKey                     string `json:"account_key,omitempty" yaml:"account_key,omitempty"`
	AccountName                    string `json:"account_name,omitempty" yaml:"account_name,omitempty"`
	ContainerName                  string `json:"container_name,omitempty" yaml:"container_name,omitempty"`
	JsonKey                        string `json:"json_key,omitempty" yaml:"json_key,omitempty"`
	Topic                          string `json:"topic,omitempty" yaml:"topic,omitempty"`
	ProjectId                      string `json:"project_id,omitempty" yaml:"project_id,omitempty"`
	ConnectionString               string `json:"connection_string,omitempty" yaml:"connection_string,omitempty"`
	EventHubName                   string `json:"event_hub_name,omitempty" yaml:"event_hub_name,omitempty"`
	ApiToken                       string `json:"apitoken,omitempty" yaml:"apitoken,omitempty"`
	Urls                           string `json:"urls,omitempty" yaml:"urls,omitempty"`
}

type Forwarder struct {
	Config ForwarderConfig `json:"config" yaml:"config"`
	Id     int             `json:"id" yaml:"id"`
	Name   string          `json:"name" yaml:"name"`
	Schema string          `json:"schema" yaml:"schema"`
}

func (f Forwarder) GetTableData() map[string]string {
	result := make(map[string]string)

	result["ID"] = fmt.Sprintf("%d", f.Id)
	result["Name"] = f.Name
	result["Schema"] = f.Schema

	return result
}

func (f Forwarder) GetColumns() []string {
	return []string{
		"ID", "Name", "Schema",
	}
}