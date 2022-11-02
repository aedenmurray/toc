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

func main() {
	URL := flag.String("url", "https://github.com/fastfire/deepdarkCTI/blob/main/forum.md", "Initial URL")
	skip := flag.String("skip", "", "File of domains to be skipped, separated by newlines")
	socksHost := flag.String("shost", "127.0.0.1", "TOR SOCKS Host")
	socksPort := flag.String("sport", "9050", "TOR SOCKS Port")
	flag.Parse()

	toSkip := state.Create()
	if *skip != "" {
		file, err := os.Open(*skip)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			domain := scanner.Text()
			toSkip.Store(domain)
		}	
	}

	client, err := tor.CreateClient(*socksHost, *socksPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	waitGroup := sync.WaitGroup{}
	hooks := new(tor.Hooks)

	node := &tor.Node{
		URL:       *URL,
		Client:    client,
		WaitGroup: &waitGroup,
		Hooks:     hooks,
		Depth:     0,
	}

	node.OnRequest(func (node *tor.Node) {
		if *skip == "" {
			return
		}

		parsed, err := url.Parse(node.URL)
		if err != nil {
			return
		}

		domain := (func() string {
			parts := strings.Split(parsed.Hostname(), ".")
			return fmt.Sprintf("%s.%s", parts[len(parts) - 2], parts[len(parts) - 1])
		})()

		if toSkip.Exists(domain) {
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
