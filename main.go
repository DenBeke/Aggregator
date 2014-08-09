package main

import "log"

func main() {
	
	//fetcher, err := NewFetcher("config/config.json")
	parser, err := NewParser("config/config.json", "config/schema.json")
	
	if(err != nil) {
		log.Println(err.Error())
		return
	}
	
	parser.Run()
	
	log.Printf("Parsed %d articles in %d feeds\n", parser.CountArticles(), parser.CountFeeds() )	
	
	parser.Save()
}