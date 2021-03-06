package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const defaultPath = "../documents"

// Configuration stores setting values
type Configuration struct {
	AppName          string `json:"appName"`
	ConfigServerAddr string `json:"configServerAddr"`
}

// Config shares the global configuration
var (
	Config *Configuration
)

// LoadBootstrapConfig reads bootstrap config file
func LoadBootstrapConfig() error {
	environment := os.Getenv("ENVRM")
	if environment == "" {
		environment = "dev"
	}

	file, err := os.Open(fmt.Sprintf("config/bootstrap-%s.json", environment))
	if err != nil {
		return err
	}

	Config = new(Configuration)
	err = json.NewDecoder(file).Decode(&Config)
	if err != nil {
		return err
	}

	return nil
}

// LoadConfigForApp loads configuration from config server
// requires appName, ENVRM env, ConfigServerAddr
func LoadConfigForApp(item interface{}) error {
	appName := Config.AppName

	environment := os.Getenv("ENVRM")
	if environment == "" {
		environment = "dev"
	}

	var client http.Client
	addr := fmt.Sprintf("%s/config/%s-%s.json", Config.ConfigServerAddr, appName, environment)

	r, _ := http.NewRequest("GET", addr, nil)
	response, err := client.Do(r)

	if err != nil {
		return err
	}
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, item)

	_ = LoadSecrets(item, environment)

	if err != nil {
		return err
	}

	return nil
}

func LoadSecrets(item interface{}, environment string) error {

	testTempPath := os.Getenv("TEST_TEMP_PATH")
	secretsFile := fmt.Sprintf("secrets-%s.json", environment)

	file, err := os.Open(path.Join(testTempPath, "config", secretsFile))

	if err != nil {
		return nil
	}
	secretData, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	return json.Unmarshal(secretData, item)
}

// AbsolutePathOfFile returns absolute path for file
func AbsolutePathOfFile(userID, filePath, fileName string) string {
	relativePath := path.Join(defaultPath, userID, filePath, fileName)
	abs, err := filepath.Abs(relativePath)

	if err != nil {
		log.Error(err)
		return ""
	}

	return abs

}

// LoadPort returns port value
func LoadPort(defaultPort string) string {

	var port string

	if envPort, ok := os.LookupEnv("PORT"); ok {
		port = fmt.Sprintf(":%s", envPort)
	} else {
		port = fmt.Sprintf(":%s", defaultPort)
	}

	return port
}
