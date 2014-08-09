package main

import "github.com/PuerkitoBio/goquery"
import "log"
//import "time"


//parseFeed parses the given feed
func (p *parser) parseFeed(schema *FeedSchema, feed *feed) {
	
	p.parseIndex(schema, feed, schema.Url)
	
}


//parseIndex will parse the given url of the given feed
func (p *parser) parseIndex(schema *FeedSchema, feed *feed, url string) {
	
	
	//Make new goquery document
	document, err := p.getDocument(url)
	if err != nil {
		log.Println("Could not parse feed: " + err.Error())
	}
	
	
	//Find all link (to article) items
	links := document.Selection.Find(schema.Index.ArticleLinks)
		
	links.EachWithBreak(func(i int, s *goquery.Selection) bool {
		
		url, exists := s.Attr("href")
		if !exists {
			log.Println("Could not parse article: No href tag found")
			return false
		}
		
		p.parseArticle(schema, feed, url)
		
		return true
		
	})
	
	
	//Check for next page
	next, exists := document.Selection.Find(schema.Index.NextLink).Attr("href")
	if !exists {
		//No other pages
		log.Println("END: No more pages found")
		return
	} else {
		//Other pages found : recursive call
		log.Println("Parsing next page")
		p.parseIndex(schema, feed, next)
	}
	
}