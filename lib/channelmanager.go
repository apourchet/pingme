package ping

import (
	b64 "encoding/base64"
	"fmt"
	"io"
	"sync"
)

// TODO Syncronize all this
type ChannelManager struct {
	Channels map[string]Channel
	sync.RWMutex
}

type Channel []chan string

func NewChannelManager() ChannelManager {
	return ChannelManager{Channels: make(map[string]Channel)}
}

func (m *ChannelManager) CreateChannel(id string) {
	m.Lock()
	_, ok := m.Channels[id]
	if !ok {
		m.Channels[id] = make(Channel, 0)
	}
	m.Unlock()
}

func (m *ChannelManager) AddListener(id string, c chan string) {
	m.Lock()
	m.Channels[id] = append(m.Channels[id], c)
	m.Unlock()
}

func (m *ChannelManager) WriteData(w io.Writer, msg string) {
	msg = b64.URLEncoding.EncodeToString([]byte(msg))
	fmt.Fprintf(w, "data: %s\n\n", msg)
}

// TODO Cleanup the closed channels somehow
func (m *ChannelManager) PingChannel(id, msg string) int {
	listeners := 0
	if c, ok := m.Channels[id]; ok {
		for _, out := range c {
			select {
			case out <- msg:
				listeners += 1
			default:
			}
		}
	}
	return listeners
}
