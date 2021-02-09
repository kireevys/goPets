package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

var existsLinks sync.Map

type Page struct {
	*Page
	link     url.URL
	_content string
	_title   string
	stopWord string
	t        Transport
}

type Chain []Page

func (c Chain) asLinks() string {
	var formatted string
	for i, page := range c {
		formatted += fmt.Sprintf("%d: %s\n", i+1, page.link.String())
	}
	return formatted
}

func (c Chain) asText() string {
	var formatted string
	for _, page := range c {
		formatted += fmt.Sprintf("-> %s", page.title())
	}
	return formatted[2:]
}

func (p Page) content() string {
	if p._content == "" {
		p._content = p.loadContent()
	}
	return p._content
}

func (p Page) title() string {
	if p._title == "" {
		titleRegex := regexp.MustCompile("<title>(.*)</title>")
		p._title = titleRegex.FindStringSubmatch(p.content())[1]
	}

	return p._title
}

func NewPage(page *Page, link url.URL, stopWord string, t Transport) Page {
	var newPage = Page{
		Page:     page,
		link:     link,
		stopWord: stopWord,
		t:        t,
	}
	if !strings.HasSuffix(link.Host, "wikipedia.org") {
		log.Println("Not *.Wikipedia.org start page.")
		os.Exit(1)
	}
	return newPage
}

func (p Page) loadContent() string {
	res, err := p.t.Get(p.link.String())
	if err != nil {
		log.Print("Transport broken")
		time.Sleep(1)
		res = p.loadContent()
	}
	return res
}

func (p Page) isFile(link url.URL) bool {
	var result = false
	for _, suffix := range filesExtension {
		if strings.HasSuffix(link.String(), suffix) {
			result = true
			break
		}
	}
	return result
}

func (p Page) extractChildren() Chain {
	var result Chain
	r := regexp.MustCompile("(?m:href=\"(/wiki/\\S+)\")")
	parsed := r.FindAllStringSubmatch(p.content(), -1)

	for _, res := range parsed {
		link, err := url.Parse(fmt.Sprintf("%s://%s%s", p.link.Scheme, p.link.Host, res[1]))
		if err != nil {
			continue
		}
		if p.isFile(*link) {
			continue
		}

		if _, ok := existsLinks.Load(link); !ok {
			page := NewPage(&p, *link, p.stopWord, p.t)
			result = append(result, page)
			existsLinks.Store(link, struct{}{})
		}
	}
	p.clear()
	return result
}

func (p Page) hasStopWord() bool {
	have, err := regexp.MatchString(fmt.Sprintf("(?i)%s", p.stopWord), p.content())
	if err != nil {
		return false
	}
	return have
}

func (p Page) getChain() Chain {
	var chain Chain
	if p.Page != nil {
		chain = append(chain, p.Page.getChain()...)
	}
	return append(chain, p)
}

func (p Page) findStopWord(c chan Page) {
	log.Printf("Start search in %s", p.title())

	if p.hasStopWord() {
		c <- p
	}

	for _, page := range p.extractChildren() {
		if page.hasStopWord() {
			c <- page
			break
		}
		go page.findStopWord(c)
	}
}

func (p Page) clear() {
	p._content = ""
}
