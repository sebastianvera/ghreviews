package graph

import (
	"sync"

	"github.com/sebastianvera/ghreviews/pkg/graph/generated"
)

type hub struct {
	mu       sync.RWMutex
	channels map[chan *generated.NewReviewEvent]bool
}

func NewHub() *hub {
	return &hub{
		channels: make(map[chan *generated.NewReviewEvent]bool),
	}
}

func (h *hub) Add(channel chan *generated.NewReviewEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.channels[channel] = true
}

func (h *hub) Remove(channel chan *generated.NewReviewEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.channels, channel)
}

func (h *hub) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.channels)
}

func (h *hub) BroadcastNewReviewEvent(event *generated.NewReviewEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for channel := range h.channels {
		channel <- event
	}
}
