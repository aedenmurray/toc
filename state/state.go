package state

import (
	"sync"
)

type State struct {
	visitedURLs map[string]struct{}
	sync.RWMutex
}

func (state *State) UpdateVisited(key string) {
	state.Lock()
	defer state.Unlock()

	state.visitedURLs[key] = struct{}{}
}

func (state *State) HasVisited(key string) bool {
	state.RLock()
	defer state.RUnlock()

	_, exists := state.visitedURLs[key]
	return exists
}

func Create() *State {
	return &State{
		visitedURLs: make(map[string]struct{}),
	}
}
