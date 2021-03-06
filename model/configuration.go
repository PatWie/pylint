// PyLint-GO
// Patrick Wieschollek <mail@patwie.com>

package model

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	GitHub struct {
		IntegrationID int    `yaml:"integration_id"`
		Secret        string `yaml:"secret"`
	} `yaml:"github"`
	Pylint struct {
		Name         string `yaml:"name"`
		Port         int    `yaml:"port"`
		URL          string `yaml:"url"`
		ReportsPath  string `yaml:"reports_path"`
		KeyFile      string `yaml:"key_file"`
		AdminId      int64  `yaml:"admin_id"`
		DatabaseFile string `yaml:"database_file"`
	} `yaml:"pylint"`
	Redis struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"redis"`
}

// Print configuration to console
func (c *Configuration) Debug() {
	fmt.Printf("IntegrationID: %v\n", c.GitHub.IntegrationID)
	fmt.Printf("ReportsPath: %v\n", c.Pylint.ReportsPath)
	fmt.Printf("Url: %v\n", c.Pylint.URL)
	fmt.Printf("Name: %v\n", c.Pylint.Name)
	fmt.Printf("Port: %v\n", c.Pylint.Port)
	fmt.Printf("DatabaseFile: %v\n", c.Pylint.DatabaseFile)
	fmt.Printf("AdminId: %v\n", c.Pylint.AdminId)
}

// parse configuration from environment
func GetConfiguration() (conf *Configuration) {
	config_fn := os.Getenv("PYLINT_CONFIGURATION")
	if config_fn == "" {
		fmt.Println("no config (env-var PYLINT_CONFIGURATION is empty)")
		return &Configuration{}
	}

	y, err := ioutil.ReadFile(config_fn)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(y, &conf)
	if err != nil {
		panic(err)
	}

	// secret could be in env-var PYLINT_GITHUB_SECRET
	secret := os.Getenv("PYLINT_GITHUB_SECRET")
	if secret != "" {
		// overwrite secret
		conf.GitHub.Secret = secret
	}

	return conf
}
