package pcap

import (
	"github.com/ff14wed/xivnet/v3"
	"github.com/thejerf/suture"
	"go.uber.org/zap"
)

type ffxivStreamFactory struct {
	cfg AdapterConfig

	providerPool *ProviderPool
	streamSupervisor *suture.Supervisor

	logger *zap.Logger
}

// New sets up the stream decoder
func (fsf *ffxivStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	// TODO: Figure out how to check if flow is outbound
	var isOutbound bool


	portStr := transport.Dst().String()
	if isOutbound {
		portStr := transport.Src().String()
	}

	lp, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		fsf.logger.Warn("Unable to parse port number", zap.String("port", portStr))
		return nil
	}
	localPort := uint16(lp)

	provider := fsf.providerPool.Get(localPort)

	frames := provider.inboundFrames
	if isOutbound {
		frames = provider.outboundFrames
	}

	streamLogName := fmt.Sprintf("stream-in-%d", localPort)
	if isOutbound {
		streamLogName := fmt.Sprintf("stream-out-%d", localPort)
	}

	// Create a new stream when a new flow is detected
	fsf.cfg.StreamUp <- provider

	// TODO: How do we figure out when the flow is gone?

	st := &ffxivStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
		frames:    frames,
		logger:    fsf.logger.Named(streamLogName),
		localPort: localPort,
		cleanup: cleanup,
	}

	fsf.streamSupervisor.Add(st)

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &stream.r
}

// ffxivStream will handle the actual decoding of FFXIV packets
type ffxivStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
	frames         chan *xivnet.Frame
	logger         *zap.Logger
	localPort      uint16
}

func (fs *ffxivStream) Serve() {
	fs.logger.Info("Running stream processor",
		zap.Object("net", fs.net),
		zap.Object("transport", fs.transport),
	)

	d := xivnet.NewDecoder(&fs.r, 32768)
	defer fs.r.Close()
	for {
		decodedFrame, err := d.NextFrame()
		// We must read until we see an EOF... very important!
		if err == nil {
			decodedFrame.CorrectTimestamps(decodedFrame.Time)
			fs.frames <- decodedFrame
			continue
		}

		fs.logger.Error(err,
			zap.Object("net", fs.net),
			zap.Object("transport", fs.transport),
		)
		switch err.(type) {
		case xivnet.EOFError:
			fs.logger.Info("Closing stream...")
			return
		case xivnet.InvalidHeaderError:
			fs.logger.Debug("Discarding bytes until we reach a good state...")
			d.DiscardDataUntilValid()
			fs.logger.Debug("Header looks correct now. Continuing...")
		}
	}
}

