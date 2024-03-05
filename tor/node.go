package tor

import (
	"io"
	"sync"

	"github.com/aedenmurray/toc/parse"
	"github.com/aedenmurray/toc/state"
)

var visited = state.Create()

type OnRequest func(node *Node)
type OnResponse func(node *Node)

type Hooks struct {
	OnRequest  OnRequest
	OnResponse OnResponse
}

type Node struct {
	URL       string
	Ref       string
	Body      []byte
	Client    *Client
	WaitGroup *sync.WaitGroup
	Hooks     *Hooks
	Depth     uint
	Skip      bool
}

func (node *Node) OnRequest(hook OnRequest) {
	node.Hooks.OnRequest = hook
}

func (node *Node) OnResponse(hook OnResponse) {
	node.Hooks.OnResponse = hook
}

func (node *Node) Crawl() {
	if visited.Exists(node.URL) {
		return
	}

	visited.Store(node.URL)
	node.WaitGroup.Add(1)

	go func() {
		defer node.WaitGroup.Done()

		if node.Hooks.OnRequest != nil {
			node.Hooks.OnRequest(node)
		}

		if node.Skip {
			return
		}

		response, err := node.Client.Get(node.URL)
		if err != nil {
			return
		}

		defer response.Body.Close()
		node.Body, err = io.ReadAll(response.Body)
		if err != nil {
			return
		}

		links := make(chan string)
		go parse.Links(node.URL, &node.Body, links)

		for link := range links {
			childNode := &Node{
				URL:       link,
				Ref:       node.URL,
				Client:    node.Client,
				WaitGroup: node.WaitGroup,
				Hooks:     node.Hooks,
				Depth:     node.Depth + 1,
			}

			childNode.Crawl()
		}

		if node.Hooks.OnResponse != nil {
			node.Hooks.OnResponse(node)
		}
	}()
}
