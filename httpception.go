package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"

	"httpception/frontend"
)

var listenAddress string
var sendAddress string

func init() {
	flag.StringVar(&listenAddress, "listen", "", "Address to listen for new connections (ex: localhost:3333)")
	flag.StringVar(&sendAddress, "send", "", "Address to listen for new connections (ex: localhost:4444)")
}

func showHelpAndExit(message string) {
	fmt.Printf("\nError: %s\n", message)
	flag.Usage()
	os.Exit(1)
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
	frontend := frontend.Frontend(frontend.NewWebSocketFrontEnd(updateChan, commandChan))
	go frontend.Start()

	// handle incoming connections
	handler := NewHttpRequestHandler(connectionChannel, errorChan, frontend.InterceptRequest, frontend.InterceptResponse)
	go handler.Start()
	<-quit
}

type HttpRequestHandler struct {
	connectionChannel <-chan net.Conn
	errorChan         chan<- error

	// interceptors
	interceptRequest  func(*http.Request) *http.Request
	interceptResponse func(*http.Response) *http.Response
}

func NewHttpRequestHandler(connectionChannel <-chan net.Conn, errorChan chan<- error, interceptRequest func(*http.Request) *http.Request, interceptResponse func(*http.Response) *http.Response) *HttpRequestHandler {
	return &HttpRequestHandler{
		connectionChannel: connectionChannel,
		errorChan:         errorChan,
		interceptRequest:  interceptRequest,
		interceptResponse: interceptResponse,
	}
}

func (h *HttpRequestHandler) Start() {
	for {
		func() {
			conn := <-h.connectionChannel
			defer conn.Close()

			// read/parse request
			reader := bufio.NewReader(conn)
			req, err := http.ReadRequest(reader)
			if err != nil {
				h.errorChan <- fmt.Errorf("Failed to parse http request: %s", err)
			}
			if req != nil {

				// intercept the request
				req = h.interceptRequest(req)

				// forward the request
				response, err := h.forwardRequest(req)
				if err != nil {
					h.errorChan <- err
				}

				// intercept the response
				response = h.interceptResponse(response)

				// send back the response to the caller
				if response != nil {
					if err := response.Write(conn); err != nil {
						h.errorChan <- errors.New(fmt.Sprintf("Failed to write response: %s", err))
					}
				}
			}
		}()
	}
}

func (h *HttpRequestHandler) forwardRequest(request *http.Request) (*http.Response, error) {

	// forward request
	conn, err := net.Dial("tcp", sendAddress)
	if err != nil {
		fmt.Errorf("Failed to dial: %s", sendAddress)
	}
	request.Write(conn)
	reader := bufio.NewReader(conn)
	return http.ReadResponse(reader, request)
}
