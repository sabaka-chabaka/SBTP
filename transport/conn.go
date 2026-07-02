package transport

import (
	"SBTP/frame"
	"bufio"
	"net"
	"time"
)

const (
	DefaultReadTimeout  = 30 * time.Second
	DefaultWriteTimeout = 30 * time.Second
)

type Conn struct {
	conn         net.Conn
	reader       *bufio.Reader
	writer       *bufio.Writer
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewConn(c net.Conn) *Conn {
	return &Conn{
		conn:         c,
		reader:       bufio.NewReader(c),
		writer:       bufio.NewWriter(c),
		readTimeout:  DefaultReadTimeout,
		writeTimeout: DefaultWriteTimeout,
	}
}

func (c *Conn) SetReadTimeout(d time.Duration) {
	c.readTimeout = d
}

func (c *Conn) SetWriteTimeout(d time.Duration) {
	c.writeTimeout = d
}

func (c *Conn) ReadFrame() (*frame.Frame, error) {
	if c.readTimeout > 0 {
		if err := c.conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
			return nil, err
		}
	}
	return frame.ReadFrame(c.reader)
}

func (c *Conn) WriteFrame(f *frame.Frame) error {
	if c.writeTimeout > 0 {
		if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return err
		}
	}
	if err := frame.WriteFrame(c.writer, f); err != nil {
		return err
	}
	return c.writer.Flush()
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Conn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
