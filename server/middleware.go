package server

import (
	"SBTP/frame"
	"log"
)

type Middleware func(Handler) Handler

func Chain(h Handler, middlewares ...Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func Logging(logger *log.Logger) Middleware {
	return func(next Handler) Handler {
		return func(req *Request) *Response {
			resp := next(req)
			logger.Printf("%s %s -> %d", req.Method, req.Path, resp.Status)
			return resp
		}
	}
}

func Recover() Middleware {
	return func(next Handler) Handler {
		return func(req *Request) (resp *Response) {
			defer func() {
				if r := recover(); r != nil {
					resp = NewResponse(frame.StatusInternalError, nil)
				}
			}()
			return next(req)
		}
	}
}
