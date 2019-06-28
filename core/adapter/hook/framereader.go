package hook

import (
	"encoding/hex"
	"fmt"

	"github.com/ff14wed/aetherometer/core/message"
	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	"go.uber.org/zap"
)

// FrameReader reads in envelopes from the hook, parses them, and converts
// the data into xivnet.Frames
type FrameReader struct {
	streamID            uint32
	envelopesChan       <-chan Envelope
	frameDecoderFactory message.DecoderFactory
	logger              *zap.Logger

	ingressFramesChan chan *xivnet.Frame
	egressFramesChan  chan *xivnet.Frame

	stop     chan struct{}
	stopDone chan struct{}
}

// NewFrameReader creates a new FrameReader provided a data source and a
// factory to create a new frame decoder
func NewFrameReader(
	streamID uint32,
	envelopesChan <-chan Envelope,
	frameDecoderFactory message.DecoderFactory,
	logger *zap.Logger,
) *FrameReader {
	return &FrameReader{
		streamID:            streamID,
		envelopesChan:       envelopesChan,
		frameDecoderFactory: frameDecoderFactory,
		logger:              logger.Named("frame-reader"),

		ingressFramesChan: make(chan *xivnet.Frame, 5000),
		egressFramesChan:  make(chan *xivnet.Frame, 5000),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the frame reader. It is responsible for sorting incoming envelopes
// of data to the correct decoder and forwarding all the decoded frames to
// subscribers.
func (d *FrameReader) Serve() {
	defer close(d.stopDone)

	ingressMuxDecoder := message.NewMuxDecoder(d.frameDecoderFactory)
	egressMuxDecoder := message.NewMuxDecoder(d.frameDecoderFactory)

	d.logger.Info("Running")
	for {
		select {
		case e, ok := <-d.envelopesChan:
			if !ok {
				continue
			}
			switch e.Op {
			case OpDebug:
				d.logger.Debug("Hook Message", zap.Uint32("data", e.Data), zap.ByteString("additional", e.Additional))
			case OpRecv:
				d.feedDataAndSendBlocks(e.Data, e.Additional, ingressMuxDecoder, false)
			case OpSend:
				d.feedDataAndSendBlocks(e.Data, e.Additional, egressMuxDecoder, true)
			default:
			}
		case <-d.stop:
			d.logger.Info("Stopping...")
			return
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning.
func (d *FrameReader) Stop() {
	close(d.stop)
	<-d.stopDone
}

// SubscribeIngress provides a channel on which consumers can listen for
// processed ingress frames decoded from the envelopes.
func (d *FrameReader) SubscribeIngress() <-chan *xivnet.Frame {
	return d.ingressFramesChan
}

// SubscribeIngress provides a channel on which consumers can listen for
// processed egress frames decoded from the envelopes.
func (d *FrameReader) SubscribeEgress() <-chan *xivnet.Frame {
	return d.egressFramesChan
}

func (d *FrameReader) feedDataAndSendBlocks(
	ctx uint32,
	data []byte,
	muxDecoder *message.MuxDecoder,
	isEgress bool,
) {
	muxDecoder.WriteData(ctx, data)
	for {
		frame, err := muxDecoder.NextFrame()
		if err != nil {
			d.logger.Error("Error reading next frame",
				zap.Uint32("context", ctx),
				zap.String("data", hex.EncodeToString(data)),
				zap.Bool("isEgress", isEgress),
				zap.Error(err),
			)
			return
		}
		if frame == nil {
			return
		}
		frame.CorrectTimestamps(frame.Time)
		blocks := frame.Blocks
		if isEgress {
			blocks = message.ProcessBlocks(frame)
		}
		var parsedBlocks []*xivnet.Block
		for _, b := range blocks {
			parsedBlock, err := datatypes.ParseBlock(b, isEgress)
			if err != nil {
				d.logger.Error("Error unmarshaling block",
					zap.Bool("isEgress", isEgress),
					zap.Uint16("opcode", b.IPCHeader.Opcode),
					zap.String("data", fmt.Sprintf("%#v", b.Data)),
					zap.Error(err),
				)
				parsedBlocks = append(parsedBlocks, b)
			} else {
				parsedBlocks = append(parsedBlocks, parsedBlock)
			}
		}
		frame.Blocks = parsedBlocks
		if isEgress {
			d.egressFramesChan <- frame
		} else {
			d.ingressFramesChan <- frame
		}
	}
}
