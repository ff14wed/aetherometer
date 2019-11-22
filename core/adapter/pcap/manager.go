package pcap

import (
	"github.com/thejerf/suture"
)

var (
	snapshotLen int32 = 8192
	promiscuous       = false
	timeout           = pcap.BlockForever
)

type assemblerFlusher struct {
	assembler *tcpassembly.Assembler
	done <-chan struct{}
}

func (a *assemblerFlusher) Serve() {
	ticker := time.Tick(1 * time.Second)
	s.assembler.FlushWithOptions(tcpassembly.FlushOptions{
		T: time.Now(),
	})
	for {
		select {
		case <-ticker:
			a.assembler.FlushWithOptions(tcpassembly.FlushOptions{
				T: time.Now().Add(-1 * time.Second),
			})
		case <-done:
			return
		}
	}
}

func (a *assemblerFlusher) Stop() {
	close(a.done)
}

type packetForwarder struct {
	source *gopacket.PacketSource
	packetChan chan gopacket.Packet
}

func (p *packetForwarder) Serve() {
	for {
		packet, err := p.source.NextPacket()
		if err == io.EOF {
			return
		} else if err == nil {
			p.packetChan <- packet
		}
	}
}

type packetCapturer struct {
	assembler *tcpassembly.Assembler
	packetChan chan gopacket.Packet
	done <-chan struct{}
}

func (p *packetCapturer) Serve() {
	for {
		select {
		case <-done:
			return nil
		case packet := <-p.packetChan:
			if packet == nil {
				continue
			}
			if packet.NetworkLayer() == nil ||
				packet.TransportLayer() == nil ||
				packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				// Unusable packet
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			p.assembler.AssembleWithTimestamp(
				packet.NetworkLayer().NetworkFlow(),
				tcp,
				packet.Metadata().Timestamp,
			)
		}
	}
}

type Manager struct {
	*suture.Supervisor
}

// TODO: Returning an error here feels yucky
func NewManager(streamPool *tcpAssembly.StreamPool, devices []string) (*Manager, error) {
	assembler := tcpassembly.NewAssembler(streamPool)
	packetChan := make(chan gopacket.Packet, 2000)

	m := &Manager{}

	af := &assemblerFlusher{
		assembler: assembler,
		done: make(chan struct{}),
	}

	m.Add(af)

	// Listen on ipv4 devices
	for i, dev := range s.devices {
		// Could probably do something with this handle
		handle, err := pcap.OpenLive(dev, snapshotLen, promiscuous, timeout)
		if err != nil {
			return err
		}
		defer func() {
			go handle.Close()
		}()

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		pf := &packetForwarder{
			source: packetSource,
			packetChan: packetChan,
		}

		m.Add(pf)
	}

	pc := &packetCapturer{
		done: make(chan struct{}),
	}
	m.Add(pc)

	return m
}
