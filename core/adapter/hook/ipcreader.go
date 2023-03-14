package hook

import (
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/ff14wed/xivnet/v3"
	"github.com/ff14wed/xivnet/v3/datatypes"
	"go.uber.org/zap"
)

// IPCReader reads in payloads from the hook, parses them, and converts
// the data into xivnet.Blocks
type IPCReader struct {
	streamID     uint32
	payloadsChan <-chan Payload
	logger       *zap.Logger

	zoneUsesPDK bool
	pdk         byte

	ingressBlocksChan chan *xivnet.Block
	egressBlocksChan  chan *xivnet.Block

	lastEgressMoveTimestamp time.Time

	stop     chan struct{}
	stopDone chan struct{}
}

// NewIPCReader creates a new IPCReader provided a data source and a
// factory to create a new IPC decoder
func NewIPCReader(
	streamID uint32,
	payloadsChan <-chan Payload,
	logger *zap.Logger,
) *IPCReader {
	return &IPCReader{
		streamID:     streamID,
		payloadsChan: payloadsChan,
		logger:       logger.Named("ipc-reader"),

		ingressBlocksChan: make(chan *xivnet.Block, 5000),
		egressBlocksChan:  make(chan *xivnet.Block, 5000),

		stop:     make(chan struct{}),
		stopDone: make(chan struct{}),
	}
}

// Serve runs the block reader. It is responsible for sorting incoming payloads
// of data to the correct decoder and forwarding all the decoded frames to
// subscribers.
func (d *IPCReader) Serve() {
	defer close(d.stopDone)

	d.logger.Info("Running")
	for {
		select {
		case e, ok := <-d.payloadsChan:
			if !ok {
				continue
			}
			switch e.Op {
			case OpDebug:
				d.logger.Debug("Hook Message", zap.Uint32("channel", e.Channel), zap.ByteString("data", e.Data))
			case OpRecv:
				d.decodeDataAndSendBlock(e.Channel, e.Data, false)
			case OpSend:
				d.decodeDataAndSendBlock(e.Channel, e.Data, true)
			case OpPing:
			case OpExit:
			default:
			}
		case <-d.stop:
			d.logger.Info("Stopping...")
			return
		}
	}
}

// Stop will shutdown this service and wait on it to stop before returning.
func (d *IPCReader) Stop() {
	close(d.stop)
	<-d.stopDone
}

// SubscribeIngress provides a channel on which consumers can listen for
// processed ingress frames decoded from the payloads.
func (d *IPCReader) SubscribeIngress() <-chan *xivnet.Block {
	return d.ingressBlocksChan
}

// SubscribeIngress provides a channel on which consumers can listen for
// processed egress frames decoded from the payloads.
func (d *IPCReader) SubscribeEgress() <-chan *xivnet.Block {
	return d.egressBlocksChan
}

func (d *IPCReader) decodeDataAndSendBlock(
	channel uint32,
	data []byte,
	isEgress bool,
) {
	if len(data) < 48 {
		d.logger.Error("Error decoding payload: not enough data: expected at least 32 bytes",
			zap.Bool("isEgress", isEgress),
			zap.String("data", hex.EncodeToString(data)),
		)
		return
	}
	block := xivnet.Block{
		SubjectID: binary.LittleEndian.Uint32(data[0:4]),
		CurrentID: binary.LittleEndian.Uint32(data[4:8]),
		Type:      xivnet.BlockTypeIPC,
		IPCHeader: xivnet.IPCHeader{
			Reserved: binary.LittleEndian.Uint16(data[16:18]),
			Opcode:   binary.LittleEndian.Uint16(data[18:20]),
			Pad2:     binary.LittleEndian.Uint16(data[20:22]),
			ServerID: binary.LittleEndian.Uint16(data[22:24]),
			Pad3:     binary.LittleEndian.Uint32(data[28:32]),
		},
	}
	timestamp := binary.LittleEndian.Uint64(data[8:16])
	msecSinceEpoch := time.Duration(timestamp) * time.Millisecond
	block.Time = time.Unix(0, 0).Add(msecSinceEpoch)

	var blockData xivnet.GenericBlockData = make([]byte, len(data)-32)

	var bd xivnet.BlockData

	if channel == 1 {
		bd = datatypes.NewBlockData(block.Opcode, isEgress)
	} else if channel == 2 {
		bd = datatypes.NewChatBlockData(block.Opcode, isEgress)
	}

	if bd == nil {
		bd = blockData
	}

	err := datatypes.UnmarshalBlockBytes(data[32:], bd)
	if err != nil {
		d.logger.Error("Error unmarshaling block:",
			zap.Bool("isEgress", isEgress),
			zap.String("data", hex.EncodeToString(data)),
			zap.Error(err),
		)
		return
	}

	block.Data = bd

	d.runPacketHooks(&block)

	if isEgress {
		switch block.Data.(type) {
		case *datatypes.EgressMovement:
			if block.Time != d.lastEgressMoveTimestamp {
				d.lastEgressMoveTimestamp = block.Time
				d.egressBlocksChan <- &block
			}
		case *datatypes.EgressInstanceMovement:
			if block.Time != d.lastEgressMoveTimestamp {
				d.lastEgressMoveTimestamp = block.Time
				d.egressBlocksChan <- &block
			}
		default:
			d.egressBlocksChan <- &block
		}

	} else {
		d.ingressBlocksChan <- &block
	}
}

func (d *IPCReader) runPacketHooks(parsedBlock *xivnet.Block) {
	if z, ok := parsedBlock.Data.(*datatypes.InitZone); ok {
		d.pdk = 0
		switch z.TerritoryTypeID {
		case 946, 947, 948, 949:
		default:
			d.zoneUsesPDK = false
			return
		}
		d.zoneUsesPDK = true
	}

	if d.zoneUsesPDK {
		if d.pdk == 0 {
			if z, ok := parsedBlock.Data.(*datatypes.Casting); ok {
				d.pdk = byte(z.ActionID - uint32(z.ActionIDName))
				d.logger.Info("Detected PDK for current zone", zap.Uint8("pdk", d.pdk))
			}
		}
		if z, ok := parsedBlock.Data.(*datatypes.Control); ok {
			if z.Type == 0x22 && d.pdk > 0 {
				z.P1 = z.P1 - uint32(d.pdk)
				return
			}
		}
	}
}
