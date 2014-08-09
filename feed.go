package main


import "sync"
import "io/ioutil"
import "encoding/json"


type feed struct {

	Title		string
	Description	string
	Articles		[]Article
	
	items 		sync.RWMutex
	

}


//NewFeed constructor
func NewFeed(title string, description string) *feed {
	
	f := new(feed)
	f.Title = title
	f.Description = description
	
	f.items = sync.RWMutex{} // ------> Is this really needed? <------
	
	return f
	
}



//AddArticle adds a new article to a feed
func (feed *feed) AddArticle(article Article) {
	
	feed.items.Lock()
	defer feed.items.Unlock()
	
	feed.Articles = append(feed.Articles, article)
	
}



//CountArticles will count all articles in the feed
func (feed *feed) CountArticles() uint {
	
	feed.items.RLock()
	defer feed.items.RUnlock()
	
	return uint( len(feed.Articles) )
	
}


//Save exports the feed to a JSON file
func (feed *feed) Save(file string) error {
	
	feed.items.RLock()
	defer feed.items.RUnlock()
	
	json, err := json.MarshalIndent(feed, "", "    ") 
	if err != nil {
		return err
	}
	
	return ioutil.WriteFile(file, json, 0777)
	
}

