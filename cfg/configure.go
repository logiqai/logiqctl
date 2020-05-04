package cfg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
)

/*
This help to configure the logiq-ctl
*/

const (
	CONFIG_DIR  = ".logiqctl"
	CONFIG_FILE = "config.toml"
	CONFIG_DB   = "logiqctl.db"
)

func Configure() (*Profiles, error) {
	ROOT, err := homedir.Dir()
	if err != nil {
		fmt.Print("Cannot get user home directory")
		return nil, err
	}
	configDir := path.Join(ROOT, CONFIG_DIR)
	exists, err := exists(configDir)
	if err != nil {
		fmt.Print("Cannot get config directory")
		return nil, err
	}
	if !exists {
		os.MkdirAll(configDir, os.ModePerm)
	}
	configFilePath := path.Join(configDir, CONFIG_FILE)

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		config, err := getNewConfig()
		if err != nil {
			return nil, err
		}
		profile := &Profiles{Configs: []Config{*config}}
		createConfig(configFilePath, profile)
		return profile, nil
	}
	return reConfigure()
}

func reConfigure() (*Profiles, error) {
	profiles, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	whatPrompt := promptui.Select{
		Label: "Config Exists, What do you want to do?",
		Items: []string{
			"List available profiles",
			"Add a new profile",
		},
	}
	what, _, err := whatPrompt.Run()
	if err != nil {
		return nil, err
	}
	switch what {
	case 0:
		PrintConfig(profiles)
		break
	case 1:
		newConfig, err := getNewConfig()
		if err != nil {
			return nil, err
		}
		profiles.Configs = append(profiles.Configs, *newConfig)
		createConfig(GetConfigFilePath(), profiles)
	}
	return profiles, nil
}

func createConfig(fileName string, profiles *Profiles) {
	buf := new(bytes.Buffer)
	toml.NewEncoder(buf).Encode(profiles)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	fmt.Printf("Created following profile at %s\n", fileName)
	fmt.Println("====================================")
	fmt.Println(buf.String())
	fmt.Println("====================================")
}

func getNewConfig() (*Config, error) {

	clusterPrompt := promptui.Prompt{
		Label: "Enter the Logiq cluster URL: ",
		//Validate: validateFile,
		Default: "logiq.mydomain.com",
	}
	cluster, err := clusterPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil, err
	}

	defaultConfig := true

	return &Config{Name: cluster, Cluster: cluster + ":8081", ApiKey: "N/A", Default: defaultConfig}, nil
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err != nil, err
}
