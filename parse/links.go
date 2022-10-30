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

func Links(body io.Reader, channel chan<- string) {
	tokenizer := html.NewTokenizer(body)
	defer close(channel)

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

				isOnionURL := OnionURLWithSchemeReg.MatchString(attribute.Val)
				if !isOnionURL {
					continue
				}

				parsedURL, err := url.Parse(attribute.Val)
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
