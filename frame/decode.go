package frame

import (
	"encoding/binary"
	"io"
)

func ReadFrame(r io.Reader) (*Frame, error) {
	header := make([]byte, HeaderSize)
	if _, err := io.ReadFull(r, header); err != nil {
		return nil, err
	}

	magic := binary.BigEndian.Uint32(header[0:4])
	if magic != Magic {
		return nil, ErrInvalidMagic
	}

	f := &Frame{
		Version: header[4],
		Type:    FrameType(header[5]),
		Flags:   Flags(binary.BigEndian.Uint16(header[6:8])),
		Status:  binary.BigEndian.Uint16(header[8:10]),
	}

	if !f.Type.Valid() {
		return nil, ErrMalformedFrame
	}

	f.MetadataLen = binary.BigEndian.Uint32(header[12:16])
	f.PayloadLen = binary.BigEndian.Uint64(header[16:24])
	copy(f.Checksum[:], header[24:56])

	if f.MetadataLen > MaxMetadataSize {
		return nil, ErrMetadataTooLarge
	}
	if f.PayloadLen > MaxPayloadSize {
		return nil, ErrFrameTooLarge
	}

	if f.MetadataLen > 0 {
		metaBytes := make([]byte, f.MetadataLen)
		if _, err := io.ReadFull(r, metaBytes); err != nil {
			return nil, err
		}
		headers, err := decodeMetadata(metaBytes)
		if err != nil {
			return nil, err
		}
		f.Metadata = headers
	}

	if f.PayloadLen > 0 {
		f.Payload = make([]byte, f.PayloadLen)
		if _, err := io.ReadFull(r, f.Payload); err != nil {
			return nil, err
		}
	}

	if err := f.VerifyChecksum(); err != nil {
		return nil, err
	}

	return f, nil
}
