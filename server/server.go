package server

import (
	"SBTP/frame"
	"SBTP/transport"
	"log"
	"net"
)

type Server struct {
	router     *Router
	middleware []Middleware
	logger     *log.Logger
}

func New() *Server {
	return &Server{
		router: NewRouter(),
		logger: log.Default(),
	}
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
	conn := transport.NewConn(rawConn)
	defer conn.Close()

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
