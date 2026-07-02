package main

import (
	"SBTP/client"
	"fmt"
	"log"
	"time"
)

func main() {
	c := client.New("localhost:9000", client.WithTimeout(5*time.Second))

	ping(c)
	echo(c)
	getUser(c)
	slowWithShortTimeout()
}

func ping(c *client.Client) {
	req := client.NewRequest("GET", "/ping", nil)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("ping failed: %v", err)
	}

	fmt.Printf("ping -> status=%s body=%s\n", resp.Status, resp.Payload)
}

func echo(c *client.Client) {
	req := client.NewRequest("POST", "/echo", []byte("hello sbtp")).
		WithChecksum()

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("echo failed: %v", err)
	}

	contentType, _ := resp.GetHeader("content-type")
	fmt.Printf("echo -> status=%s content-type=%s body=%s\n", resp.Status, contentType, resp.Payload)
}

func getUser(c *client.Client) {
	req := client.NewRequest("GET", "/users/42", nil)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalf("getUser failed: %v", err)
	}

	if resp.Status.IsSuccess() {
		fmt.Printf("getUser -> %s\n", resp.Payload)
	} else {
		fmt.Printf("getUser failed with status %s\n", resp.Status)
	}
}

func slowWithShortTimeout() {
	c := client.New("localhost:9000", client.WithTimeout(1*time.Second))
	req := client.NewRequest("GET", "/slow", nil)

	start := time.Now()
	_, err := c.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("slow -> timed out after %v as expected (%v)\n", elapsed, err)
	} else {
		fmt.Println("slow -> unexpectedly succeeded, timeout did not trigger")
	}
}
