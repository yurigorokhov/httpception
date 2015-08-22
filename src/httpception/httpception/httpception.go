package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"httpception/frontend"
)

var listenAddress string
var sendAddress string
var debuggingAddress string

func init() {
	flag.StringVar(&listenAddress, "listen", "", "Address to listen for new connections (ex: localhost:3333)")
	flag.StringVar(&sendAddress, "send", "", "Address to listen for new connections (ex: localhost:4444)")
	flag.StringVar(&debuggingAddress, "debug", ":9999", "Address to listen for debugging connection (default: :9999)")
}

func main() {
	flag.Parse()
	if len(listenAddress) == 0 {
		showHelpAndExit("listen is a required parameter")
	}
	if len(sendAddress) == 0 {
		showHelpAndExit("send is a required parameter")
	}

	// start listening for connections
	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Printf("Error listening on %s: %s", listenAddress, err)
		os.Exit(1)
	}

	// accept connections indefinitely
	connectionChannel := make(chan net.Conn)
	errorChan := make(chan error)
	quit := make(chan struct{})
	go func() {
		for {

			//TODO: handle accept errors
			conn, _ := l.Accept()
			connectionChannel <- conn
		}
	}()

	// handle errors
	go func() {
		for {
			select {
			case err := <-errorChan:
				fmt.Errorf("ERROR: %s", err)
			}
		}
	}()

	// initialize frontend
	updateChan := make(chan frontend.Update)
	commandChan := make(chan frontend.Command)
	frontend := frontend.Frontend(frontend.NewWebSocketFrontend(updateChan, commandChan, debuggingAddress))
	go frontend.Start()

	// handle incoming connections
	handler := NewHTTPProxy(connectionChannel, errorChan, frontend.InterceptRequest, frontend.InterceptResponse)
	go handler.Start()
	<-quit
}

func showHelpAndExit(message string) {
	fmt.Printf("\nError: %s\n", message)
	flag.Usage()
	os.Exit(1)
}
