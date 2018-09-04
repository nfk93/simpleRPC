package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

// Task for use with RPC
type Task int

// Void type
type Void struct{}

var a int

// Increment increments a and returns a
func (t *Task) Increment(_ Void, reply *int) error {
	a++
	*reply = a
	return nil
}

func runClient() {
	var err error
	var reply int

	// Create a TCP connection to localhost on port 1234
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	for {
		client.Call("Task.Increment", struct{}{}, &reply)
		fmt.Printf("Called Task.Increment remotely and received: %d\n", reply)
		time.Sleep(2 * time.Second)
	}
}

func runServer() {
	task := new(Task)
	// Publish the receivers methods
	err := rpc.Register(task)
	if err != nil {
		log.Fatal("Format of service Task isn't correct. ", err)
	}
	// Register a HTTP handler
	rpc.HandleHTTP()
	// Listen to TPC connections on port 1234
	listener, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error: ", e)
	}
	log.Printf("Serving RPC server on port %d", 1234)
	// Start accept incoming HTTP connections
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("Error serving: ", err)
	}
}

func main() {
	a = 0

	fmt.Print("Enter 'c' for client or 's' for server... ")
	var input string
	fmt.Scanln(&input)

	if input == "c" {
		fmt.Println("Starting client...")
		runClient()
	} else if input == "s" {
		fmt.Println("Starting server...")
		runServer()
	} else {
		log.Fatal("please enter 's' or 'c'")
	}
}
