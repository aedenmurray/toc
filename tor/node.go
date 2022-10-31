package tor

import (
	"bytes"
	"io"
	"sync"

	"github.com/aedenmurray/toc/parse"
	"github.com/aedenmurray/toc/state"
)

var visitedURLs = state.Create()

type OnRequest func(node *Node)
type OnResponse func(node *Node)

type Hooks struct {
	OnRequest  OnRequest
	OnResponse OnResponse
}

type Node struct {
	URL       string
	Ref       string
	Client    *Client
	Buffer    *bytes.Buffer
	WaitGroup *sync.WaitGroup
	Hooks     *Hooks
	Depth     uint
}

func (node *Node) OnRequest(hook OnRequest) {
	node.Hooks.OnRequest = hook
}

func (node *Node) OnResponse(hook OnResponse) {
	node.Hooks.OnResponse = hook
}

func (node *Node) Crawl() {
	hasVisited := visitedURLs.Exists(node.URL)
	if hasVisited {
		return
	}

	visitedURLs.Store(node.URL)
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

		node.Buffer = new(bytes.Buffer)
		teeReader := io.TeeReader(response.Body, node.Buffer)

		linksChannel := make(chan string)
		go parse.Links(teeReader, linksChannel)
		for link := range linksChannel {
			childDepth := node.Depth + 1
			childNode := &Node{
				URL:       link,
				Ref:       node.URL,
				Client:    node.Client,
				WaitGroup: node.WaitGroup,
				Hooks:     node.Hooks,
				Depth:     childDepth,
			}

			childNode.Crawl()
		}

		if node.Hooks.OnResponse != nil {
			node.Hooks.OnResponse(node)
		}
	}()
}
