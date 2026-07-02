package frame

func (t FrameType) String() string {
	switch t {
	case TypeRequest:
		return "REQUEST"
	case TypeResponse:
		return "RESPONSE"
	case TypePing:
		return "PING"
	case TypePong:
		return "PONG"
	case TypeError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func (t FrameType) Valid() bool {
	switch t {
	case TypeRequest, TypeResponse, TypePing, TypePong, TypeError:
		return true
	default:
		return false
	}
}

func (f Flags) Has(flag Flags) bool {
	return f&flag != 0
}

func (f Flags) Set(flag Flags) Flags {
	return f | flag
}

func (f Flags) Clear(flag Flags) Flags {
	return f &^ flag
}
