package frontend

import (
	"fmt"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type WebSocketConn struct {
	Conn     *websocket.Conn
	DoneChan chan<- bool
}

type SocketHandler struct {
	connLock    *sync.Mutex
	connections []*WebSocketConn
	commandChan chan<- Command
	updateChan  <-chan Update
}

func NewSocketHandler(commandChan chan<- Command, updateChan <-chan Update) *SocketHandler {
	return &SocketHandler{
		connLock:    &sync.Mutex{},
		connections: make([]*WebSocketConn, 0, 1),
		commandChan: commandChan,
		updateChan:  updateChan,
	}
}

func (s *SocketHandler) Run() {
	for {
		update := <-s.updateChan
		if len(s.connections) > 0 {
			s.connLock.Lock()
			for i, conn := range s.connections {
				if conn == nil {
					continue
				}
				if err := websocket.JSON.Send(conn.Conn, update); err != nil {
					fmt.Printf("\nERROR] %v", err)
					conn.DoneChan <- true
					s.connections[i] = nil
				}
			}
			s.connLock.Unlock()
		}
	}
}

func (s *SocketHandler) HandleConn(ws *websocket.Conn) {
	defer ws.Close()
	doneChan := make(chan bool, 1)
	s.connLock.Lock()
	s.connections = append(s.connections, &WebSocketConn{
		Conn:     ws,
		DoneChan: doneChan,
	})
	s.connLock.Unlock()

	// read commands
forloop:
	for {
		select {
		case <-doneChan:
			break forloop
		default:
			var msg Command
			if err := websocket.JSON.Receive(ws, &msg); err != nil {
				if err == io.EOF {
					break forloop
				}
				fmt.Printf("\n[ERROR] %v", err)
			} else {
				s.commandChan <- msg
			}
		}
	}
	fmt.Printf("Done listening on socket")
}
