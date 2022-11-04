package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/aedenmurray/toc/parse"
	"github.com/aedenmurray/toc/state"
	"github.com/aedenmurray/toc/tor"
)

var domainsToSkip *state.State = state.Create()

func main() {
	URL := flag.String("url", "https://github.com/fastfire/deepdarkCTI/blob/main/forum.md", "Initial URL")
	skip := flag.String("skip", "", "File of domains to be skipped, separated by newlines")
	socksHost := flag.String("shost", "127.0.0.1", "TOR SOCKS Host")
	socksPort := flag.String("sport", "9050", "TOR SOCKS Port")
	flag.Parse()

	(func() {
		if *skip == "" {
			return
		}
	
		file, err := os.Open(*skip)
		if err != nil {
			log.Fatal(err)
		}
	
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			domain := scanner.Text()
			domainName := strings.Split(domain, ".")[0]
			domainsToSkip.Store(domainName)
		}
	})()

	hooks := new(tor.Hooks)
	waitGroup := sync.WaitGroup{}

	client, err := tor.CreateClient(*socksHost, *socksPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	node := &tor.Node{
		URL:       *URL,
		Client:    client,
		WaitGroup: &waitGroup,
		Hooks:     hooks,
		Depth:     0,
	}

	node.OnRequest(func(node *tor.Node) {
		if *skip == "" {
			return
		}

		parsed, err := url.Parse(node.URL)
		if err != nil {
			return
		}

		splitDomain := strings.Split(parsed.Hostname(), ".")
		domainName := splitDomain[len(splitDomain)-2]

		if domainsToSkip.Exists(domainName) {
			node.Skip = true
		}
	})

	node.OnResponse(func(node *tor.Node) {
		title := "(" + parse.Title(node.Buffer) + ")"
		fmt.Println(node.Depth, node.URL, title)
	})

	node.Crawl()
	waitGroup.Wait()
}
