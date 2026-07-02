package server

import (
	"SBTP/frame"
	"strings"
)

type route struct {
	segments []string
	handler  Handler
}

type Router struct {
	routes []route
}

func NewRouter() *Router {
	return &Router{}
}

func (rt *Router) Handle(path string, h Handler) {
	rt.routes = append(rt.routes, route{
		segments: splitPath(path),
		handler:  h,
	})
}

func splitPath(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return []string{}
	}
	return strings.Split(trimmed, "/")
}

func (rt *Router) dispatch(req *Request) *Response {
	reqSegments := splitPath(req.Path)

	for _, r := range rt.routes {
		params, ok := match(r.segments, reqSegments)
		if !ok {
			continue
		}
		req.params = params
		return r.handler(req)
	}

	return NewResponse(frame.StatusNotFound, nil)
}

func match(pattern, path []string) (map[string]string, bool) {
	if len(pattern) != len(path) {
		return nil, false
	}

	params := make(map[string]string)

	for i, seg := range pattern {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			name := seg[1 : len(seg)-1]
			params[name] = path[i]
			continue
		}
		if seg != path[i] {
			return nil, false
		}
	}

	return params, true
}
