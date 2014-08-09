package main

import "encoding/json"
import "os"
import "errors"


type Config struct {
	MaxConnections	uint16
	MaxAttempts 		uint16
	WaitTime 		uint
	CacheDir 		string
	SaveDir 			string
}


//readConfig will create and parse the given config file
//and return the corresponding Config object
func readConfig(file_name string) (*Config, error) {
	
	//Open the config file
	config_file, err := os.Open(file_name)
	
	if err != nil {
		return nil, errors.New("Could not open config file: " + err.Error())
	}
	
	
	//Parse the config file
	config := Config{}

	jsonParser := json.NewDecoder(config_file)
	err = jsonParser.Decode(&config)
	
	if err != nil {
		return nil, errors.New("Could not parse config file: " +  err.Error())
	} else {
		return &config, nil
	}
	
}