package main

import (
	"SBTP/frame"
	"SBTP/server"
	"fmt"
	"log"
	"os"
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

	srv.Handle("/dog", func(req *server.Request) *server.Response {
		filePath := "sabaka.png"

		data, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error:", err)
			return server.NewResponse(frame.StatusNotFound, nil)
		}

		return server.NewResponse(frame.StatusOK, data)
	})

	srv.Handle("/hello/{name}", func(req *server.Request) *server.Response {
		name := req.Param("name")
		return server.NewResponse(frame.StatusOK, []byte("Hello "+name+"!"))
	})

	log.Println("SBTP server listening on :9000")
	if err := srv.ListenAndServe(":9000"); err != nil {
		log.Fatal(err)
	}
}
