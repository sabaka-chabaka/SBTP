package frame

import (
	"bytes"
	"encoding/binary"
)

func encodeMetadata(headers []Header) ([]byte, error) {
	var buf bytes.Buffer

	for _, h := range headers {
		if len(h.Key) > 0xFFFF {
			return nil, ErrMalformedFrame
		}
		if len(h.Value) > 0xFFFF {
			return nil, ErrMalformedFrame
		}

		if err := binary.Write(&buf, binary.BigEndian, uint16(len(h.Key))); err != nil {
			return nil, err
		}
		buf.WriteString(h.Key)

		if err := binary.Write(&buf, binary.BigEndian, uint16(len(h.Value))); err != nil {
			return nil, err
		}
		buf.WriteString(h.Value)
	}

	return buf.Bytes(), nil
}

func decodeMetadata(data []byte) ([]Header, error) {
	var headers []Header
	offset := 0

	for offset < len(data) {
		if offset+2 > len(data) {
			return nil, ErrMalformedFrame
		}
		keyLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += 2

		if offset+keyLen > len(data) {
			return nil, ErrMalformedFrame
		}
		key := string(data[offset : offset+keyLen])
		offset += keyLen

		if offset+2 > len(data) {
			return nil, ErrMalformedFrame
		}
		valueLen := int(binary.BigEndian.Uint16(data[offset : offset+2]))
		offset += 2

		if offset+valueLen > len(data) {
			return nil, ErrMalformedFrame
		}
		value := string(data[offset : offset+valueLen])
		offset += valueLen

		headers = append(headers, Header{Key: key, Value: value})
	}

	return headers, nil
}

func (f *Frame) GetHeader(key string) (string, bool) {
	for _, h := range f.Metadata {
		if h.Key == key {
			return h.Value, true
		}
	}
	return "", false
}

func (f *Frame) SetHeader(key, value string) {
	for i, h := range f.Metadata {
		if h.Key == key {
			f.Metadata[i].Value = value
			return
		}
	}
	f.Metadata = append(f.Metadata, Header{Key: key, Value: value})
}
