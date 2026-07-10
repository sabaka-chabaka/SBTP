package client

import (
	"SBTP/crypto"
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
	pool    *transport.Pool
	useTLS  bool
}

type Option func(*Client)

func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}

func WithEncryption() Option {
	return func(c *Client) {
		c.useTLS = true
	}
}

func New(addr string, opts ...Option) *Client {
	c := &Client{
		addr:    addr,
		timeout: transport.DefaultReadTimeout,
		pool:    transport.NewPool(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) dial() (*transport.Conn, error) {
	rawConn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return nil, err
	}

	conn := transport.NewConn(rawConn)

	if c.useTLS {
		session, err := crypto.ClientHandshake(rawConn)
		if err != nil {
			rawConn.Close()
			return nil, err
		}
		conn.EnableEncryption(session)
	}

	return conn, nil
}

func (c *Client) Do(req *Request) (*Response, error) {
	conn, err := c.pool.Get(c.addr)
	if err != nil {
		return nil, err
	}

	conn.SetReadTimeout(c.timeout)
	conn.SetWriteTimeout(c.timeout)

	if err := conn.WriteFrame(req.toFrame()); err != nil {
		c.pool.Discard(conn)
		return nil, err
	}

	f, err := conn.ReadFrame()
	if err != nil {
		c.pool.Discard(conn)
		return nil, err
	}

	if f.Type != frame.TypeResponse {
		c.pool.Discard(conn)
		return nil, ErrUnexpectedFrameType
	}

	c.pool.Put(c.addr, conn)
	return newResponse(f), nil
}

func (c *Client) Close() {
	c.pool.CloseIdle()
}
