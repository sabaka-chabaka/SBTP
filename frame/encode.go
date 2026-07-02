package frame

import (
	"encoding/binary"
	"io"
)

func WriteFrame(w io.Writer, f *Frame) error {
	if !f.Type.Valid() {
		return ErrMalformedFrame
	}

	if len(f.Payload) > MaxPayloadSize {
		return ErrFrameTooLarge
	}

	metaBytes, err := encodeMetadata(f.Metadata)
	if err != nil {
		return err
	}

	if len(metaBytes) > MaxMetadataSize {
		return ErrMetadataTooLarge
	}

	f.MetadataLen = uint32(len(metaBytes))
	f.PayloadLen = uint64(len(f.Payload))

	header := make([]byte, HeaderSize)

	binary.BigEndian.PutUint32(header[0:4], Magic)
	header[4] = f.Version
	header[5] = uint8(f.Type)
	binary.BigEndian.PutUint16(header[6:8], uint16(f.Flags))
	binary.BigEndian.PutUint16(header[8:10], uint16(f.Status))
	binary.BigEndian.PutUint32(header[12:16], f.MetadataLen)
	binary.BigEndian.PutUint64(header[16:24], f.PayloadLen)
	copy(header[24:56], f.Checksum[:])

	if _, err := w.Write(header); err != nil {
		return err
	}

	if len(metaBytes) > 0 {
		if _, err := w.Write(metaBytes); err != nil {
			return err
		}
	}

	if len(f.Payload) > 0 {
		if _, err := w.Write(f.Payload); err != nil {
			return err
		}
	}

	return nil
}
