package main


import "log"
import "time"
import "strings"
import "github.com/kennygrant/sanitize"


//Removes starting/ending spaces from a string
func removeSpaces(input string) string {

	for input[0] == ' ' {
		input = strings.TrimPrefix(input, " ")
	}
	for input[len(input)-1] == ' ' {
		input = strings.TrimSuffix(input, " ")
	}
	
	return input
	
}



func (p *parser) parseArticle(schema *FeedSchema, feed *feed, url string) {

	document, err := p.getDocument(url)
	if err != nil {
		log.Println("Could not parse article (" + url + "): " + err.Error())
		return
	}

	//Parse title
	title, err := document.Selection.Find(schema.Article.Title).Html()
	if err != nil {
		log.Println("Could not parse article: Could not find title of article")
		return
	}

	//Parse content
	content, err := document.Selection.Find(schema.Article.Content).Html()
	if err != nil {
		log.Println("Could not parse article: Could not find content of article")
		return
	}

	//Parse date
	raw_date, err := document.Selection.Find(schema.Article.Date).Html()
	if err != nil {
		log.Println("Could not parse article: Could not find date of article")
		return
	}

	//Convert date to UNIX timestamp
	date, err := time.Parse(schema.Article.DateFormat, raw_date)
	if err != nil {
		log.Println("Could not parse article: Could parse date of article")
		return
	}
	timestamp := int(date.Unix())

	
	//Create new article
	article := Article{}
	article.Title = title
	article.Content = removeSpaces( sanitize.HTML(content) )
	article.Date = uint(timestamp)

	//Add article to feed
	feed.AddArticle(article)


	//Some logs...
	log.Println(title + " - " + raw_date)

}