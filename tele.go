package main

import (
  "fmt"
  "sync"
  "strings"
  "strconv"
  "net/http"

  "github.com/antchfx/xpath"
  "github.com/antchfx/xquery/html"
)

const (
  MAX_WORKERS = 30 // Maximum worker goroutines
  HOLDING_CAPACITY = 100 // Holding capacity of the channel
)

type Scrapper struct {
  Url string
}

func (s *Scrapper) Scrap() {
  resp, _ := http.Get(s.Url)

  defer resp.Body.Close()

  doc, _ := htmlquery.Parse(resp.Body)
  expr := xpath.MustCompile("//div[contains(@class, 'contest-item-status-disputed')]/a/text()")
  iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)

  for iter.MoveNext() {
    if i, err := strconv.Atoi(strings.TrimSpace(iter.Current().Value()[:2])); err == nil {
      issueCount += i
    }
  }

  sitesParsed += 1
}

var issueCount int = 0
var sitesParsed int = 0

func main () {
  var urls []string

  resp, _ := http.Get("https://instantview.telegram.org/contest")

  defer resp.Body.Close()

  doc, _ := htmlquery.Parse(resp.Body)
  expr := xpath.MustCompile("//div[@class='contest-item-domain']/a/@href")
  iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)

  for iter.MoveNext() {
    urls = append(urls, "https://instantview.telegram.org" + iter.Current().Value())
  }

  fmt.Printf("Sites in contest: %d \n\n", len(urls))

  urlsToScrap := make(chan *Scrapper, HOLDING_CAPACITY)

  var wg sync.WaitGroup
  for i := 0; i < MAX_WORKERS; i++ {
    wg.Add(1)
    go func() {
      for url := range urlsToScrap {
        url.Scrap()
      }
      wg.Done()
    }()
  }

  for i := 0; i < len(urls); i++ {
    urlsToScrap <- &Scrapper{Url: urls[i]}
  }
  close(urlsToScrap)

  wg.Wait()

  fmt.Printf("Parsing is done.\n\nSites parsed: %d\nCount of issues: %d\n", sitesParsed, issueCount)
}