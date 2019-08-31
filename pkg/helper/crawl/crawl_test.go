package crawl

import (
	"strings"
	"testing"
)

func TestCrawl(t *testing.T) {
	ch := make(chan string)
	go func() {
		if err := CrawlHref(BackfillBase, ch); err != nil {
			t.Fatal(err)
		}
	}()
	for u := range ch {
		if strings.HasSuffix(u, ".php") {
			t.Log(u)
		}
	}
}
