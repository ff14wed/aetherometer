package pcap

import (
	"errors"
	"sync"

	"github.com/ff14wed/xivnet/v3"
)

type ProviderPool struct {
	providers map[uint16]Provider

	lock sync.Mutex
}

func NewProviderPool() *ProviderPool {
	return &ProviderPool{
		providers: make(map[uint16]Provider),
	}
}

func (pp *ProviderPool) Get(localPort uint16) *Provider {
	pp.lock.Lock()
	defer pp.lock.Unlock()

	p, found := pp.providers[localPort]
	if !found {
		p = Provider{
			inboundFrames:  make(chan *xivnet.Frame),
			outboundFrames: make(chan *xivnet.Frame),
			localPort:      localPort,
		}
		pp.providers[localPort] = p
	}
	return p
}

type Provider struct {
	inboundFrames  chan *xivnet.Frame
	outboundFrames chan *xivnet.Frame

	localPort uint16
}

func (p *Provider) StreamID() int {
	return int(p.localPort)
}

func (p *Provider) SubscribeIngress() <-chan *xivnet.Frame {
	return p.inboundFrames
}

func (p *Provider) SubscribeEgress() <-chan *xivnet.Frame {
	return p.outboundFrames
}

func (p *Provider) SendRequest([]byte) ([]byte, error) {
	return nil, errors.New("Not supported")
}
