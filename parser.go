package main



import "github.com/PuerkitoBio/goquery"
import "strings"
import "log"
import "github.com/kennygrant/sanitize"


type parser struct {
	
	feeds 			[](*feed)
	feeds_config 	[]FeedSchema
	config			Config
	fetcher			fetcher
	
}


//NewParser constructor
func NewParser(config_file string, feeds_schema_file string) (*parser, error) {
	
	
	//Read config data
	config, err := readConfig(config_file)
	if(err != nil) {
		return nil,err
	}
	
	fetcher,err := NewFetcher(config)
	if err != nil {
		return nil,err
	}
	
	schemas, err := readFeedsConfig(feeds_schema_file)
	if err != nil {
		return nil,err
	}
	
	
	p := new(parser)
	p.fetcher = *fetcher
	p.feeds_config = *schemas
	p.config = *config
	
	return p,nil
	
}



//Run the parser
func (p *parser) Run() {
	
	for _,schema := range p.feeds_config {
		
		feed := NewFeed(schema.Title, schema.Description)
		p.feeds = append(p.feeds, feed)
		
		log.Println("BEGIN: Parsing " + schema.Title)
		p.parseFeed(&schema, feed)
		log.Println("END: Parsing " + schema.Title)
		
	}
	
}



//GetFeeds returns all parsed feeds
//this function must of-course be called after the 'Run' function
func (p *parser) GetFeeds() ( *[](*feed) ) {
	
	return &p.feeds
	
}



//Save all feeds
func (p *parser) Save() {
	
	
	for _,feed := range p.feeds {
		if nil != feed.Save( p.config.SaveDir + "/" + sanitize.Name(feed.Title) + ".json" ) {
			log.Println("Could not save feed '" + feed.Title + "'")
		}
	}
	
}



//getDocument retrieves the url from the fetcher
//the fetcher will cache the url content
func (p *parser) getDocument(url string) (*goquery.Document, error) {
	
	
	//Fetch the url
	content, err := p.fetcher.Get(url)
	if err != nil {
		return nil, err
	}
	
	//Create GoQuery document
	reader := strings.NewReader(content)
	
	document, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, err
	}
	
	return document, err
	
}



//CountFeeds will count all parsed feeds
func (p *parser) CountFeeds() uint {
	return uint( len(p.feeds) )
}


//CountArticles counts all the articles of all the feeds
func (p *parser) CountArticles() uint {
	
	var count uint = 0
	
	for _,feed := range p.feeds {
		count += feed.CountArticles()
	}
	
	return count
	
}
