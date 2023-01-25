package parse

import (
	"bytes"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var TabRegexp = regexp.MustCompile(`\t`)
var NewLineRegexp = regexp.MustCompile(`\r?\n`)

func Title(body *[]byte) string {
	reader := bytes.NewReader(*body)
	tokenizer := html.NewTokenizer(reader)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "title" {
				tokenType := tokenizer.Next()
				if tokenType == html.TextToken {
					token := tokenizer.Token()
					title := strings.TrimSpace(token.Data)
					title = TabRegexp.ReplaceAllString(title, "")
					title = NewLineRegexp.ReplaceAllString(title, "")

					return title
				}
			}
		}
	}

	return ""
}
