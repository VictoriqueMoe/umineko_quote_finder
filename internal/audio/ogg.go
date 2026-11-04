package audio

import (
	"encoding/binary"
	"fmt"
)

var oggCRCTable [256]uint32

func init() {
	for i := range 256 {
		r := uint32(i) << 24
		for range 8 {
			if r&0x80000000 != 0 {
				r = (r << 1) ^ 0x04c11db7
			} else {
				r <<= 1
			}
		}
		oggCRCTable[i] = r
	}
}

type oggPage struct {
	headerType     byte
	granulePos     int64
	serialNumber   uint32
	sequenceNumber uint32
	segmentTable   []byte
	data           []byte
}

func parseOggPages(data []byte) ([]oggPage, error) {
	var pages []oggPage
	offset := 0

	for offset < len(data) {
		if offset+27 > len(data) {
			return nil, fmt.Errorf("truncated page header at offset %d", offset)
		}
		if string(data[offset:offset+4]) != "OggS" {
			return nil, fmt.Errorf("missing OggS capture pattern at offset %d", offset)
		}

		headerType := data[offset+5]
		granulePos := int64(binary.LittleEndian.Uint64(data[offset+6 : offset+14]))
		serialNumber := binary.LittleEndian.Uint32(data[offset+14 : offset+18])
		sequenceNumber := binary.LittleEndian.Uint32(data[offset+18 : offset+22])
		numSegments := int(data[offset+26])

		segTableEnd := offset + 27 + numSegments
		if segTableEnd > len(data) {
			return nil, fmt.Errorf("truncated segment table at offset %d", offset)
		}

		segmentTable := make([]byte, numSegments)
		copy(segmentTable, data[offset+27:segTableEnd])

		var dataSize int
		for _, s := range segmentTable {
			dataSize += int(s)
		}

		dataEnd := segTableEnd + dataSize
		if dataEnd > len(data) {
			return nil, fmt.Errorf("truncated page data at offset %d", offset)
		}

		pageData := make([]byte, dataSize)
		copy(pageData, data[segTableEnd:dataEnd])

		pages = append(pages, oggPage{
			headerType:     headerType,
			granulePos:     granulePos,
			serialNumber:   serialNumber,
			sequenceNumber: sequenceNumber,
			segmentTable:   segmentTable,
			data:           pageData,
		})

		offset = dataEnd
	}

	return pages, nil
}

func (p *oggPage) serialize() []byte {
	headerSize := 27 + len(p.segmentTable)
	buf := make([]byte, headerSize+len(p.data))

	copy(buf[0:4], "OggS")
	buf[4] = 0
	buf[5] = p.headerType
	binary.LittleEndian.PutUint64(buf[6:14], uint64(p.granulePos))
	binary.LittleEndian.PutUint32(buf[14:18], p.serialNumber)
	binary.LittleEndian.PutUint32(buf[18:22], p.sequenceNumber)
	binary.LittleEndian.PutUint32(buf[22:26], 0)
	buf[26] = byte(len(p.segmentTable))
	copy(buf[27:headerSize], p.segmentTable)
	copy(buf[headerSize:], p.data)

	crc := p.oggCRC(buf)
	binary.LittleEndian.PutUint32(buf[22:26], crc)

	return buf
}

func (*oggPage) oggCRC(data []byte) uint32 {
	var crc uint32
	for _, b := range data {
		crc = (crc << 8) ^ oggCRCTable[(crc>>24)^uint32(b)]
	}
	return crc
}
