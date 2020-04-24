package cfg

import (
	"fmt"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/prologic/bitcask"
)

var CONFIG *Config

func LoadConfig() (*Profiles, error) {

	configFileName := GetConfigFilePath()

	profiles := &Profiles{}
	_, err := toml.DecodeFile(configFileName, profiles)
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

func PrintConfig(profiles *Profiles) {
	for _, p := range profiles.Configs {
		fmt.Printf("\nProfile Name: %s\n", p.Name)
		fmt.Printf("Cluster: %s\n", p.Cluster)
		fmt.Printf("Api Key: %s\n", p.ApiKey)
		fmt.Printf("Default Profile: %v\n", p.Default)
	}
}

func GetConfigFilePath() string {
	ROOT, err := getRootFolder()
	if err != nil {
		fmt.Print("Cannot get user home directory")
	}
	return path.Join(ROOT, CONFIG_DIR, CONFIG_FILE)
}

func getRootFolder() (string, error) {
	ROOT, err := homedir.Dir()
	return ROOT, err
}

func GetConfigDB() *bitcask.Bitcask {
	ROOT, err := getRootFolder()
	if err != nil {
		fmt.Print("Cannot get user home directory")
	}
	dbPath := path.Join(ROOT, CONFIG_DIR, CONFIG_DB)
	db, _ := bitcask.Open(dbPath)
	return db
}
