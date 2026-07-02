package frame

import "errors"

const (
	Magic uint32 = 0x53425450
)

const (
	HeaderSize   = 56
	ChecksumSize = 32
)

type FrameType uint8

const (
	TypeRequest FrameType = iota
	TypeResponse
	TypePing
	TypePong
	TypeError
)

type Flags uint16

const (
	FlagNone     Flags = 0
	FlagChecksum Flags = 1 << 0
)

type Header struct {
	Key   string
	Value string
}

type Frame struct {
	Version     uint8
	Type        FrameType
	Flags       Flags
	Status      uint16
	MetadataLen uint32
	PayloadLen  uint64
	Checksum    [ChecksumSize]byte
	Metadata    []Header
	Payload     []byte
}

var (
	ErrInvalidMagic       = errors.New("frame: invalid magic")
	ErrUnsupportedVersion = errors.New("frame: unsupported version")
	ErrFrameTooLarge      = errors.New("frame: payload exceeds max size")
	ErrMetadataTooLarge   = errors.New("frame: metadata exceeds max size")
	ErrChecksumMismatch   = errors.New("frame: checksum mismatch")
	ErrMalformedFrame     = errors.New("frame: malformed frame")
)
