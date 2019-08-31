package crawl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	BackfillBase = "https://www.football-data.co.uk/data.php"
)

func CrawlHref(u string, ch chan string) error {
	defer close(ch)
	parsedURL, err := url.Parse(u)
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocument(u)
	if err != nil {
		return err
	}
	doc.Find("a[href]").Each(func(i int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		fmt.Println(href)
		switch {
		case strings.HasPrefix(href, "http"):
			ch <- href
		case strings.HasPrefix(href, "/"):
			ch <- fmt.Sprintf("%s://%s%s", parsedURL.Scheme, parsedURL.Host, href)
		case len(href) != 0:
			ch <- fmt.Sprintf("%s://%s/%s", parsedURL.Scheme, parsedURL.Host, href)
		}
	})
	return nil
}
