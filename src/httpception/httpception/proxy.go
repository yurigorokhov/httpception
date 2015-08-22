package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// HTTPProxy proxies requests while allowing them to be intercepted
type HTTPProxy struct {
	connectionChannel <-chan net.Conn
	errorChan         chan<- error
	sendAddress       string

	// interceptors
	interceptRequest  func(*http.Request) *http.Request
	interceptResponse func(*http.Response) *http.Response
}

// NewHTTPProxy creates a new proxy
func NewHTTPProxy(
	connectionChannel <-chan net.Conn,
	errorChan chan<- error,
	interceptRequest func(*http.Request) *http.Request,
	interceptResponse func(*http.Response) *http.Response,
	sendAddress string) *HTTPProxy {
	return &HTTPProxy{
		connectionChannel: connectionChannel,
		errorChan:         errorChan,
		interceptRequest:  interceptRequest,
		interceptResponse: interceptResponse,
		sendAddress:       sendAddress,
	}
}

// Start starts the proxy
func (h *HTTPProxy) Start() {
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

				// rewrite the request Host header
				req = h.rewriteRequest(req)

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

func (h *HTTPProxy) forwardRequest(request *http.Request) (*http.Response, error) {

	// forward request
	conn, err := net.Dial("tcp", h.sendAddress)
	if err != nil {
		fmt.Errorf("Failed to dial: %s", h.sendAddress)
	}
	request.Write(conn)
	reader := bufio.NewReader(conn)
	return http.ReadResponse(reader, request)
}

func (h *HTTPProxy) rewriteRequest(request *http.Request) *http.Request {
	request.Host = h.sendAddress
	return request
}
