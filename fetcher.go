package main

import "net/http"
import "io/ioutil"
import "errors"
import "time"
import "crypto/md5"
import "encoding/hex"
import "os"
import "sync"




type fetcher struct {
	cache sync.RWMutex
	sem chan int
	config Config
}



//NewFetcher creates a new 'fetcher' object
//and initializes it with a configuration
func NewFetcher(config *Config) (*fetcher, error) {
	
	fetcher := new(fetcher)
	
	//Save values
	fetcher.config = *config
	fetcher.sem = make(chan int, fetcher.config.MaxConnections)
	
	return fetcher,nil
	
}




//Get the content of the url
//The function will first try to get the url from the cache
//and will (if needed) fetch it from the remote host
//
//Configuration settings can be done in the `config/config.json` file
func (fetcher *fetcher) Get(url string) (string, error) {

	
	//Check if url can be found on disk
	content, err := fetcher.getDiskUrl(url)
	if err == nil {
		
		//url found on disk, return content
		return content, nil
	
	} else {
	
		content, err = fetcher.getRemoteUrl(url)
		
		if err != nil {
			return "",err
		} else {
			_ = fetcher.writeDiskUrl(url, content)
			return content,nil
		}
		
	}
	
	//Retrieve url from remote host
	
	return "", nil
}



//getRemoteUrl will fetch the content of a remote url
//
//Config values for max attempts, wait time and max connections
//can be set in the `config/config.json` file
func (fetcher *fetcher) getRemoteUrl(url string) (string, error) {
	
	
	//Try a couple of times to fetch the url
	//this prevents errors when there are temporary issues at the remote
	
	var resp *http.Response
	var err error
	
	for i := 0; i < int(fetcher.config.MaxAttempts); i++ {
	
		//Semaphore implementation to prevent to many open connections
		// Wait for active queue to drain
		fetcher.sem <- 1
		
		//Actual http request
		resp, err = http.Get(url)
		
		// Done: enable next request to run
		<- fetcher.sem
		
		if err != nil {
			return "", errors.New("Could not fetch url")
		}
		
		if resp.StatusCode == 200 {
			//Status code 200 : everything fine
			//OK!!!
			break;
		} else if resp.StatusCode/100 == 5 && i < 9 {
			//Status code 5xx : server problem
			//Try again after a bit of time
			resp.Body.Close()
			time.Sleep(time.Second * time.Duration(fetcher.config.WaitTime))
		} else {
			//Other status code : can't fetch url
			resp.Body.Close()
			return "", errors.New("Could not fetch url")
		}
	
		
	}
	
	
	//Get the content of the URL
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return "", errors.New("Could not read response")
	} else {
		resp.Body.Close()
		return string(content),nil
	}
	
	
}



//getDiskUrl will retrieve the url from the disk cache
func (fetcher *fetcher) getDiskUrl(url string) (string, error) {
	
	fetcher.cache.RLock()
	defer fetcher.cache.RUnlock()
	
	//Create file name	
	filename := fetcher.config.CacheDir + "/" + fetcher.hash(url) + ".html"

	//Check if file exists
	if _, err := os.Stat(filename); err != nil {
		return "", errors.New("Could not retrieve local disk cache")
	}
	
	//If file exists : try to get content
	raw_content, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", errors.New("Could not retrieve local disk cache")
	}
	
	//Return string containg the file
	return string(raw_content), nil

	
}


//writeDiskUrl will save a cache entry for the given url/content
func (fetcher *fetcher) writeDiskUrl(url string, content string) error {
	
	fetcher.cache.Lock()
	defer fetcher.cache.Unlock()
	
	//Create file name	
	filename := fetcher.config.CacheDir + "/" + fetcher.hash(url) + ".html"
	
	//Write to disk
	err := ioutil.WriteFile(filename, []byte(content), 0755)
	if err != nil {
		return errors.New("Could not write file to cache")
	}
	
	return nil
	
}



//hash the given url
func (fetcher *fetcher) hash(url string) string {
	raw_hash := md5.Sum([]byte(url))
	hash := hex.EncodeToString(raw_hash[:])
	
	return hash
}

