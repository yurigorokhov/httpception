package frontend

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"sync"

	"golang.org/x/net/websocket"
)

// Frontend represents a debugging iterface
type Frontend interface {
	InterceptRequest(*http.Request) *http.Request
	InterceptResponse(*http.Response) *http.Response
	Start()
}

// WebSocketFrontend represents the main web interface
type WebSocketFrontend struct {
	updateChan       chan UpdateInterface
	commandChan      chan Command
	debuggingAddress string

	settingsMutex    *sync.Mutex
	debuggingEnabled bool
	isPaused         bool
	continueChannel  chan struct{}
}

// NewWebSocketFrontend creates a new WebSocketFrontend
func NewWebSocketFrontend(
	updateChan chan UpdateInterface,
	commandChan chan Command,
	debuggingAddress string) *WebSocketFrontend {
	return &WebSocketFrontend{
		updateChan:       updateChan,
		commandChan:      commandChan,
		debuggingAddress: debuggingAddress,
		debuggingEnabled: false,
		settingsMutex:    &sync.Mutex{},
		continueChannel:  make(chan struct{}),
	}
}

// Start starts up the frontend
func (f *WebSocketFrontend) Start() {

	// handle
	newConnectionChan := make(chan struct{})
	socketHandler := NewSocketHandler(f.commandChan, f.updateChan, newConnectionChan)
	go socketHandler.Run()

	// listen for commands
	go func() {
		for {
			select {
			case command := <-f.commandChan:
				switch command.Type {
				case ContinueCommand:
					f.settingsMutex.Lock()
					if f.isPaused {
						f.settingsMutex.Unlock()
						f.continueChannel <- struct{}{}
					} else {
						f.settingsMutex.Unlock()
					}
				case EnableDebuggingCommand:
					f.settingsMutex.Lock()
					f.debuggingEnabled = true
					f.settingsMutex.Unlock()
					f.updateChan <- NewDebuggingToggleMessage(true)
				case DisableDebuggingCommand:
					f.settingsMutex.Lock()
					f.debuggingEnabled = false
					if f.isPaused {
						go func() { f.continueChannel <- struct{}{} }()
						f.isPaused = false
					}
					f.settingsMutex.Unlock()
					f.updateChan <- NewDebuggingToggleMessage(false)
				}
			case <-newConnectionChan:
				f.settingsMutex.Lock()
				f.updateChan <- NewInitialUpdateMessage(f.debuggingEnabled)
				f.settingsMutex.Unlock()
			}
		}
	}()

	// handle websocket connections
	http.Handle("/_socket", websocket.Handler(socketHandler.HandleConn))
	http.Handle("/", http.FileServer(http.Dir("./src/httpception/frontend/web/")))
	fmt.Printf("Listening on: %s\n", f.debuggingAddress)
	http.ListenAndServe(f.debuggingAddress, nil)
}

// InterceptRequest allows the debugger to view and modify the request
func (f *WebSocketFrontend) InterceptRequest(request *http.Request) *http.Request {
	b, _ := httputil.DumpRequest(request, true)
	f.updateChan <- NewRequestUpdateMessage(string(b))

	// only wait for debugger command if debugging is turned on
	if f.debuggingEnabled {
		f.settingsMutex.Lock()
		f.isPaused = true
		f.settingsMutex.Unlock()
		<-f.continueChannel
		f.settingsMutex.Lock()
		f.isPaused = false
		f.settingsMutex.Unlock()
	}
	return request
}

// InterceptResponse allows the debugger to view and modify the response
func (f *WebSocketFrontend) InterceptResponse(response *http.Response) *http.Response {
	b, _ := httputil.DumpResponse(response, true)
	f.updateChan <- NewResponseUpdateMessage(string(b))

	// only wait for debugger command if debugging is turned on
	if f.debuggingEnabled {
		f.settingsMutex.Lock()
		f.isPaused = true
		f.settingsMutex.Unlock()
		<-f.continueChannel
		f.settingsMutex.Lock()
		f.isPaused = false
		f.settingsMutex.Unlock()
	}
	return response
}

var _ = Frontend(&WebSocketFrontend{})
