package frontend

import (
	"net/http"
	"net/http/httputil"
	"sync"

	"golang.org/x/net/websocket"
)

type Frontend interface {
	InterceptRequest(*http.Request) *http.Request
	InterceptResponse(*http.Response) *http.Response
	Start()
}

type WebSocketFrontend struct {
	updateChan  chan Update
	commandChan chan Command

	settingsMutex    *sync.Mutex
	debuggingEnabled bool
	isPaused         bool
	continueChannel  chan struct{}
}

func NewWebSocketFrontEnd(updateChan chan Update, commandChan chan Command) *WebSocketFrontend {
	return &WebSocketFrontend{
		updateChan:       updateChan,
		commandChan:      commandChan,
		debuggingEnabled: false,
		settingsMutex:    &sync.Mutex{},
		continueChannel:  make(chan struct{}),
	}
}

func (f *WebSocketFrontend) Start() {

	// handle
	socketHandler := NewSocketHandler(f.commandChan, f.updateChan)
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
					f.updateChan <- Update{Type: DebuggingEnabledUpdate, Value: ""}
				case DisableDebuggingCommand:
					f.settingsMutex.Lock()
					f.debuggingEnabled = false
					if f.isPaused {
						go func() { f.continueChannel <- struct{}{} }()
						f.isPaused = false
					}
					f.settingsMutex.Unlock()
					f.updateChan <- Update{Type: DebuggingDisabledUpdate, Value: ""}
				}
			}
		}
	}()

	// handle websocket connections
	http.Handle("/_socket", websocket.Handler(socketHandler.HandleConn))
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":8081", nil)
}

func (f *WebSocketFrontend) InterceptRequest(request *http.Request) *http.Request {

	b, _ := httputil.DumpRequest(request, true)
	f.updateChan <- Update{Type: NewRequestUpdate, Value: string(b)}

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

func (f *WebSocketFrontend) InterceptResponse(response *http.Response) *http.Response {
	b, _ := httputil.DumpResponse(response, true)
	f.updateChan <- Update{Type: NewResponseUpdate, Value: string(b)}

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
