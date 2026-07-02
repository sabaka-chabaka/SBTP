package client

import (
	"SBTP/frame"
	"SBTP/transport"
	"errors"
	"net"
	"time"
)

var ErrUnexpectedFrameType = errors.New("client: unexpected frame type in response")

type Client struct {
	addr    string
	timeout time.Duration
}

type Option func(*Client)

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}

func New(addr string, opts ...Option) *Client {
	c := &Client{
		addr:    addr,
		timeout: transport.DefaultReadTimeout,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Do(req *Request) (*Response, error) {
	rawConn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return nil, err
	}

	conn := transport.NewConn(rawConn)
	defer conn.Close()

	conn.SetReadTimeout(c.timeout)
	conn.SetWriteTimeout(c.timeout)

	if err := conn.WriteFrame(req.toFrame()); err != nil {
		return nil, err
	}

	f, err := conn.ReadFrame()
	if err != nil {
		return nil, err
	}

	if f.Type != frame.TypeResponse {
		return nil, ErrUnexpectedFrameType
	}

	return newResponse(f), nil
}
