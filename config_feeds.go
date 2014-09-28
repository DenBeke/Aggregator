package main


import "encoding/json"
import "os"
import "errors"


type FeedSchema struct {
	
	Title 		string
	Url 			string
	Description string
	
	
	Index struct {
		ArticleLinks	string
		LinksPrefix	string
		NextLink		string
	}
	
	
	Article struct {
		Title 		string
		Date 		string
		DateFormat 	string
		Content 		string
	}
	
}


//readConfig will create and parse the given feeds config schema
//and return the corresponding FeedSchema array
func readFeedsConfig(file_name string) (*[]FeedSchema, error) {

	//Open the config file
	config_file, err := os.Open(file_name)

	if err != nil {
		return nil, errors.New("Could not open feeds schema file: " + err.Error())
	}


	//Parse the config file
	config := []FeedSchema{}

	jsonParser := json.NewDecoder(config_file)
	err = jsonParser.Decode(&config)

	if err != nil {
		return nil, errors.New("Could not parse feeds schema file: " +  err.Error())
	} else {
		return &config, nil
	}

}