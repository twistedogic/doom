package crawl

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func CrawlHref(u string, ch chan string) error {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}
	log.Printf("crawling %s", u)
	doc, err := goquery.NewDocument(u)
	if err != nil {
		return err
	}
	doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		switch {
		case strings.HasPrefix(href, "http"):
			ch <- href
		case len(href) != 0:
			ch <- fmt.Sprintf("%s://%s/%s", parsedURL.Scheme, parsedURL.Host, strings.TrimPrefix(href, "/"))
		}
	})
	return nil
}
