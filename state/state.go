package state

import (
	"sync"
)

type State struct {
	internal map[string]struct{}
	sync.RWMutex
}

func (state *State) Store(key string) {
	state.Lock()
	defer state.Unlock()

	state.internal[key] = struct{}{}
}

func (state *State) Exists(key string) bool {
	state.RLock()
	defer state.RUnlock()

	_, exists := state.internal[key]
	return exists
}

func Create() *State {
	return &State{
		internal: make(map[string]struct{}),
	}
}
