package ts

import "io"
import "bytes"
import "errors"
import "github.com/32bitkid/mpeg/util"

var (
	EOP = errors.New("end of packet")
)

type streamState int

const (
	_                 = iota
	ready streamState = iota
	drained
)

func NewPayloadUnitReader(source io.Reader, where PacketTester) io.Reader {
	return &payloadUnitBuffer{
		br:             util.NewBitReader(source),
		where:          where,
		startIndicator: where.And(IsPayloadUnitStart),
		state:          drained,
	}
}

type payloadUnitBuffer struct {
	currentPacket  *Packet
	buffer         bytes.Buffer
	br             util.BitReader32
	where          PacketTester
	startIndicator PacketTester
	state          streamState
}

func (stream *payloadUnitBuffer) Read(p []byte) (n int, err error) {
	if stream.state == drained {
		_, ferr := stream.fill()
		if ferr != nil {
			return 0, ferr
		}
		stream.state = ready
	}

	for len(p) > 0 {
		rn, rerr := stream.buffer.Read(p)
		n += rn
		p = p[rn:]

		if rerr == io.EOF {
			stream.state = drained
			return n, EOP
		} else if rerr != nil {
			return n, rerr
		}
	}

	return
}

func (stream *payloadUnitBuffer) fill() (n int, err error) {

	// initialize
	if stream.currentPacket == nil {
		stream.currentPacket = new(Packet)

		// step until first start indicator
		for {
			isStart := stream.startIndicator(stream.currentPacket)
			if isStart {
				break
			}
			err = stream.currentPacket.ReadFrom(stream.br)
			if err != nil {
				return
			}
		}
	}

	// Read until next start indicator
	for {
		if stream.where(stream.currentPacket) {
			cn, err := stream.buffer.Write(stream.currentPacket.Payload)
			n += cn
			if err != nil {
				break
			}
		}

		err = stream.currentPacket.ReadFrom(stream.br)
		if err != nil {
			break
		}

		if stream.startIndicator(stream.currentPacket) {
			break
		}
	}

	return

}