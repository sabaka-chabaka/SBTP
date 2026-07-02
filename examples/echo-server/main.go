package main

import (
	"SBTP/frame"
	"SBTP/server"
	"log"
	"time"
)

func main() {
	srv := server.New()

	srv.Use(server.Recover())
	srv.Use(server.Logging(log.Default()))

	srv.Handle("/ping", func(req *server.Request) *server.Response {
		return server.NewResponse(frame.StatusOK, []byte("pong"))
	})

	srv.Handle("/echo", func(req *server.Request) *server.Response {
		return server.NewResponse(frame.StatusOK, req.Payload).
			SetHeader("content-type", "text/plain").
			WithChecksum()
	})

	srv.Handle("/users/42", func(req *server.Request) *server.Response {
		body := []byte(`{"id":42,"name":"Sabaka"}`)
		return server.NewResponse(frame.StatusOK, body).
			SetHeader("content-type", "application/json")
	})

	srv.Handle("/slow", func(req *server.Request) *server.Response {
		time.Sleep(2 * time.Second)
		return server.NewResponse(frame.StatusOK, []byte("done"))
	})

	log.Println("SBTP server listening on :9000")
	if err := srv.ListenAndServe(":9000"); err != nil {
		log.Fatal(err)
	}
}
