package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
)

var filesExtension = []string{".svg", ".png", ".jpg", ".js"}

var (
	StopWord  = "Гитлер"
	StartPage = "https://ru.wikipedia.org"
)

func main() {
	flag.StringVar(&StopWord, "w", StopWord, "Word for search")
	flag.StringVar(&StartPage, "p", StartPage, "Start page for search")

	flag.Parse()

	log.Println(fmt.Sprintf("Word %s", StopWord))
	log.Println(fmt.Sprintf("Start page %s\n", StartPage))

	page, _ := url.Parse(StartPage)
	var p = NewPage(nil, *page, StopWord, NewHTTP())

	c := make(chan Page)
	go p.findStopWord(c)

	result := <-c

	fmt.Printf("\nPath by title\n%+v\n\n", result.getChain().asText())
	fmt.Printf("\nPath by links\n%+v\n\n", result.getChain().asLinks())
	log.Printf("Result page: %v", result.link.String())

	fmt.Printf("Request count: %d\n", result.t.count())
}
