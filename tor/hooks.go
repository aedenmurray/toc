package tor

import "bytes"

type OnRequest func(node *Node)
type OnResponse func(node *Node, buffer *bytes.Buffer)

type Hooks struct {
	OnRequest  OnRequest
	OnResponse OnResponse
}
