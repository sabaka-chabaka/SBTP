package server

import "SBTP/frame"

type Router struct {
	routes map[string]Handler
}

func NewRouter() *Router {
	return &Router{make(map[string]Handler)}
}

func (rt *Router) Handle(path string, h Handler) {
	rt.routes[path] = h
}

func (rt *Router) dispatch(req *Request) *Response {
	handler, ok := rt.routes[req.Path]
	if !ok {
		return NewResponse(frame.StatusNotFound, nil)
	}
	return handler(req)
}
