package client

import "SBTP/frame"

type Request struct {
	Method   string
	Path     string
	Headers  []frame.Header
	Payload  []byte
	checksum bool
}

type Response struct {
	Status  frame.Status
	Headers []frame.Header
	Payload []byte
	raw     *frame.Frame
}

func NewRequest(method, path string, payload []byte) *Request {
	return &Request{
		Method:  method,
		Path:    path,
		Payload: payload,
	}
}

func (r *Request) SetHeader(key, value string) *Request {
	r.Headers = append(r.Headers, frame.Header{Key: key, Value: value})
	return r
}

func (r *Request) WithChecksum() *Request {
	r.checksum = true
	return r
}

func (r *Request) toFrame() *frame.Frame {
	metadata := append([]frame.Header{
		{Key: "method", Value: r.Method},
		{Key: "path", Value: r.Path},
	}, r.Headers...)

	f := &frame.Frame{
		Version:  1,
		Type:     frame.TypeRequest,
		Metadata: metadata,
		Payload:  r.Payload,
	}
	if r.checksum {
		f.ApplyChecksum()
	}
	return f
}

func newResponse(f *frame.Frame) *Response {
	return &Response{
		Status:  frame.Status(f.Status),
		Headers: f.Metadata,
		Payload: f.Payload,
		raw:     f,
	}
}

func (r *Response) GetHeader(key string) (string, bool) {
	for _, h := range r.Headers {
		if h.Key == key {
			return h.Value, true
		}
	}
	return "", false
}
