package server

import (
	"SBTP/crypto"
	"SBTP/frame"
	"SBTP/transport"
	"log"
	"net"
)

type Server struct {
	router     *Router
	middleware []Middleware
	logger     *log.Logger
	requireTLS bool
}

func New() *Server {
	return &Server{
		router: NewRouter(),
		logger: log.Default(),
	}
}

func (s *Server) RequireEncryption() {
	s.requireTLS = true
}

func (s *Server) Handle(path string, h Handler) {
	s.router.Handle(path, h)
}

func (s *Server) Use(m Middleware) {
	s.middleware = append(s.middleware, m)
}

func (s *Server) ListenAndServe(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		rawConn, err := ln.Accept()
		if err != nil {
			s.logger.Printf("accept error: %v", err)
			continue
		}
		go s.handleConn(rawConn)
	}
}

func (s *Server) handleConn(rawConn net.Conn) {
	defer rawConn.Close()

	conn := transport.NewConn(rawConn)

	if s.requireTLS {
		session, err := crypto.ServerHandshake(rawConn)
		if err != nil {
			s.logger.Printf("handshake failed from %s: %v", rawConn.RemoteAddr(), err)
			return
		}
		conn.EnableEncryption(session)
	}

	for {
		f, err := conn.ReadFrame()
		if err != nil {
			return
		}

		if f.Type != frame.TypeRequest {
			continue
		}

		req := newRequest(f)
		handler := Chain(s.router.dispatch, s.middleware...)
		resp := handler(req)

		if err := conn.WriteFrame(resp.toFrame()); err != nil {
			return
		}
	}
}
