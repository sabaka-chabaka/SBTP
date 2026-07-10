package transport

import (
	"SBTP/crypto"
	"SBTP/frame"
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
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
	session      *crypto.Session
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

func (c *Conn) EnableEncryption(session *crypto.Session) {
	c.session = session
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

	if c.session == nil {
		return frame.ReadFrame(c.reader)
	}

	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(c.reader, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	ciphertext := make([]byte, length)
	if _, err := io.ReadFull(c.reader, ciphertext); err != nil {
		return nil, err
	}

	plaintext, err := c.session.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	return frame.ReadFrame(bytes.NewReader(plaintext))
}

func (c *Conn) WriteFrame(f *frame.Frame) error {
	if c.writeTimeout > 0 {
		if err := c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return err
		}
	}

	if c.session == nil {
		if err := frame.WriteFrame(c.writer, f); err != nil {
			return err
		}
		return c.writer.Flush()
	}

	var buf bytes.Buffer
	if err := frame.WriteFrame(&buf, f); err != nil {
		return err
	}

	ciphertext, err := c.session.Encrypt(buf.Bytes())
	if err != nil {
		return err
	}

	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(len(ciphertext)))
	if _, err := c.writer.Write(length); err != nil {
		return err
	}
	if _, err := c.writer.Write(ciphertext); err != nil {
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
