package parse

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var OnionURLWithSchemeReg *regexp.Regexp = regexp.MustCompile(`(^http://)\b[1-7a-z]{56}\.onion\b`)
var OnionURLReg *regexp.Regexp = regexp.MustCompile(`\b[1-7a-z]{56}\.onion\b`)

func IsImage(link string) bool {
	extensions := []string{".png", ".jpg", ".jpeg", ".gif"}
	for _, extension := range extensions {
		if strings.HasSuffix(link, extension) {
			return true
		}
	}

	return false
}

func Links(nodeURL string, body io.Reader, channel chan<- string) {
	defer close(channel)

	isOnionNode := OnionURLWithSchemeReg.MatchString(nodeURL)
	parsedNodeURL, err := url.Parse(nodeURL)
	if err != nil {
		return
	}

	tokenizer := html.NewTokenizer(body)
	sendToChannel := func(link string) {
		if IsImage(link) {
			return
		}

		channel <- link
	}

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

				href := attribute.Val
				isOnionURL := OnionURLWithSchemeReg.MatchString(href)
				isAbsoluteURL := strings.Contains(href, "http://") || strings.Contains(href, "https://")
				var toParse string

				if isOnionURL && isAbsoluteURL {
					toParse = href
				} else if isOnionNode && !isAbsoluteURL {
					hostWithScheme := parsedNodeURL.Scheme + "://" + parsedNodeURL.Host
					toParse = hostWithScheme + href
				} else {
					continue
				}

				parsedURL, err := url.Parse(toParse)
				if err != nil {
					continue
				}

				trimmedPath := strings.TrimSuffix(parsedURL.Path, "/")
				formattedLink := parsedURL.Scheme + "://" + parsedURL.Host + trimmedPath
				sendToChannel(formattedLink)
			}
		}

		if tokenType == html.TextToken {
			token := tokenizer.Token()
			onionURLs := OnionURLReg.FindAllString(token.Data, -1)
			for _, onionURL := range onionURLs {
				formattedLink := fmt.Sprintf("http://%s", onionURL)
				sendToChannel(formattedLink)
			}
		}
	}
}
