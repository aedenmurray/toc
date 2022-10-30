package tor

import (
	"bytes"
	"io"
	"sync"

	"github.com/aedenmurray/toc/parse"
	"github.com/aedenmurray/toc/state"
)

var visitedURLs = state.Create()

type Node struct {
	URL       string
	Client    *Client
	WaitGroup *sync.WaitGroup
	Hooks     *Hooks
	Depth     uint
}

func (node *Node) OnRequest(handler OnRequest) {
	node.Hooks.OnRequest = handler
}

func (node *Node) OnResponse(handler OnResponse) {
	node.Hooks.OnResponse = handler
}

func (node *Node) Crawl() {
	hasVisited := visitedURLs.HasVisited(node.URL)
	if hasVisited {
		return
	}

	visitedURLs.UpdateVisited(node.URL)
	node.WaitGroup.Add(1)

	go func() {
		defer node.WaitGroup.Done()

		if node.Hooks.OnRequest != nil {
			node.Hooks.OnRequest(node)
		}

		response, err := node.Client.Get(node.URL)
		if err != nil {
			return
		}

		var bodyBuffer bytes.Buffer
		teeReader := io.TeeReader(response.Body, &bodyBuffer)

		linksChannel := make(chan string)
		go parse.Links(teeReader, linksChannel)
		for link := range linksChannel {
			childDepth := node.Depth + 1
			childNode := &Node{
				URL:       link,
				Client:    node.Client,
				WaitGroup: node.WaitGroup,
				Hooks:     node.Hooks,
				Depth:     childDepth,
			}

			childNode.Crawl()
		}

		if node.Hooks.OnResponse != nil {
			node.Hooks.OnResponse(node, &bodyBuffer)
		}
	}()
}
