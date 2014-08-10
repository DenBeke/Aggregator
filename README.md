Aggregator
==========

This script is a simple aggregator for fetching web content of different websites.
Websites that need to be parsed, are added in the `config/schema.json`, with
their corresponding schema.

The aggregator will automatically fetch all the pages of each feed, and save them
in a JSON file stored in the `output` directory.

Configuration of the aggregator can be done in the `config/config.json` file.


Building
--------

### Install dependencies 

    go get github.com/PuerkitoBio/goquery
    go get github.com/kennygrant/sanitize



### Build & Run Aggregator

	go build
	./Aggregator


Acknowledgements
----------------

* [GoQuery](https://github.com/PuerkitoBio/goquery) by Martin Angers & Contributors
* [sanitize](https://github.com/kennygrant/sanitize) by Kenny Grant


Author
------

[Mathias Beke](http://denbeke.be)
