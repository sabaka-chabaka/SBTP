package server

import "SBTP/frame"

type Request struct {
	Method  string
	Path    string
	Headers []frame.Header
	Payload []byte
	params  map[string]string
	raw     *frame.Frame
}

type Response struct {
	Status   frame.Status
	Headers  []frame.Header
	Payload  []byte
	checksum bool
}

type Handler func(*Request) *Response

func newRequest(f *frame.Frame) *Request {
	method, _ := f.GetHeader("method")
	path, _ := f.GetHeader("path")

	return &Request{
		Method:  method,
		Path:    path,
		Headers: f.Metadata,
		Payload: f.Payload,
		raw:     f,
	}
}

func NewResponse(status frame.Status, payload []byte) *Response {
	return &Response{
		Status:  status,
		Payload: payload,
	}
}

func (r *Response) SetHeader(key, value string) *Response {
	r.Headers = append(r.Headers, frame.Header{Key: key, Value: value})
	return r
}

func (r *Response) WithChecksum() *Response {
	r.checksum = true
	return r
}

func (r *Response) toFrame() *frame.Frame {
	f := &frame.Frame{
		Version:  1,
		Type:     frame.TypeResponse,
		Status:   uint16(r.Status),
		Metadata: r.Headers,
		Payload:  r.Payload,
	}
	if r.checksum {
		f.ApplyChecksum()
	}
	return f
}

func (r *Request) Param(name string) string {
	return r.params[name]
}
