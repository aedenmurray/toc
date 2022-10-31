package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/aedenmurray/toc/parse"
	"github.com/aedenmurray/toc/tor"
)

var banner string = `
the onion cralwer - https://github.com/aedenmurray/toc
------------------------------------------------------
`

func main() {
	fmt.Println(banner)

	initialURL := flag.String("url", "https://github.com/fastfire/deepdarkCTI/blob/main/forum.md", "Initial URL")
	socksHost := flag.String("shost", "127.0.0.1", "TOR SOCKS Host")
	socksPort := flag.String("sport", "9050", "TOR SOCKS Port")
	flag.Parse()

	client, err := tor.CreateClient(*socksHost, *socksPort)
	if err != nil {
		fmt.Println(err)
		return
	}

	var waitGroup sync.WaitGroup
	node := &tor.Node{
		URL:       *initialURL,
		Client:    client,
		WaitGroup: &waitGroup,
		Hooks:     &tor.Hooks{},
		Depth:     0,
	}

	node.OnResponse(func(node *tor.Node) {
		title := "(" + parse.Title(node.Buffer) + ")"
		depth := fmt.Sprintf("%-4d", node.Depth)
		fmt.Println(depth, node.Ref, "->", node.URL, title)
	})

	node.Crawl()
	waitGroup.Wait()
}
