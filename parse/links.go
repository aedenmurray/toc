package parse

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var OnionURLWithSchemeReg *regexp.Regexp = regexp.MustCompile(`(^http://)\b[1-7a-z]{56}\.onion\b`)
var OnionURLReg *regexp.Regexp = regexp.MustCompile(`\b[1-7a-z]{56}\.onion\b`)

func ShouldSkipBasedOnExtension(link string) bool {
	extensions := []string{".png", ".jpg", ".jpeg", ".gif", ".mp4", ".mp3", ".bmp"}
	for _, extension := range extensions {
		if strings.HasSuffix(link, extension) {
			return true
		}
	}

	return false
}

func Links(nodeURL string, nodeBody *[]byte, links chan<- string) {
	hrefs := make(chan string)
	defer close(links);

	send := func(link string) {
		if ShouldSkipBasedOnExtension(link) {
			return
		}

		links <- strings.TrimSuffix(link, "/");
	}

	go func() {
		defer close(hrefs)
		reader := bytes.NewReader(*nodeBody);
		tokenizer := html.NewTokenizer(reader);
	
		for {
			tokenType := tokenizer.Next()
			if tokenType == html.ErrorToken {
				break
			}
	
			if tokenType == html.StartTagToken {
				token := tokenizer.Token()
				if token.Data != "a" {
					continue
				}
	
				for _, attribute := range token.Attr {
					if attribute.Key != "href" {
						continue
					}

					hrefs <- attribute.Val
				}
			}
		}
	}()

	func() {
		nodeURLParsed, err := url.Parse(nodeURL);
		if err != nil {
			return
		}

		for href := range hrefs {
			isAbsoluteRef := strings.Contains(href, "//")
	
			if isAbsoluteRef {
				if !OnionURLWithSchemeReg.MatchString(href) {
					continue
				}

				send(href)
				continue
			}
	
			if !OnionURLReg.MatchString(nodeURLParsed.Host) {
				continue
			}

			if strings.HasPrefix(href, "mailto:") {
				continue
			}

			if !strings.HasPrefix(href, "/") {
				href = "/" + href
			}
	
			send(nodeURLParsed.Scheme + "://" + nodeURLParsed.Host + href)
		}
	}()

	func() {
		onionURLs := OnionURLReg.FindAll(*nodeBody, -1);
		for _, onionURL := range onionURLs {
			formatted := fmt.Sprintf("http://%s", string(onionURL))
			send(formatted)
		}
	}()
}
