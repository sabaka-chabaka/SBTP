package transport

import (
	"net"
	"sync"
	"time"
)

const (
	DefaultMaxIdleConns    = 10
	DefaultIdleConnTimeout = 90 * time.Second
)

type Pool struct {
	mu          sync.Mutex
	idle        map[string][]*idleConn
	maxIdle     int
	idleTimeout time.Duration
	dialTimeout time.Duration
}

type idleConn struct {
	conn     *Conn
	returnAt time.Time
}

func NewPool() *Pool {
	return &Pool{
		idle:        make(map[string][]*idleConn),
		maxIdle:     DefaultMaxIdleConns,
		idleTimeout: DefaultIdleConnTimeout,
		dialTimeout: 5 * time.Second,
	}
}

func (p *Pool) SetMaxIdleConns(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.maxIdle = n
}

func (p *Pool) SetIdleTimeout(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.idleTimeout = d
}

func (p *Pool) Get(addr string) (*Conn, error) {
	p.mu.Lock()
	bucket := p.idle[addr]

	for len(bucket) > 0 {
		last := len(bucket) - 1
		candidate := bucket[last]
		bucket = bucket[:last]
		p.idle[addr] = bucket

		if time.Since(candidate.returnAt) > p.idleTimeout {
			p.mu.Unlock()
			candidate.conn.Close()
			p.mu.Lock()
			continue
		}

		p.mu.Unlock()
		return candidate.conn, nil
	}
	p.mu.Unlock()

	rawConn, err := net.DialTimeout("tcp", addr, p.dialTimeout)
	if err != nil {
		return nil, err
	}
	return NewConn(rawConn), nil
}

func (p *Pool) Put(addr string, conn *Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	bucket := p.idle[addr]
	if len(bucket) >= p.maxIdle {
		p.mu.Unlock()
		conn.Close()
		p.mu.Lock()
		return
	}

	p.idle[addr] = append(bucket, &idleConn{
		conn:     conn,
		returnAt: time.Now(),
	})
}

func (p *Pool) Discard(conn *Conn) {
	conn.Close()
}

func (p *Pool) CloseIdle() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for addr, bucket := range p.idle {
		for _, ic := range bucket {
			ic.conn.Close()
		}
		delete(p.idle, addr)
	}
}
