package frontend

// CommandType is the type of command
type CommandType uint

const (

	// EnableDebuggingCommand is a command from the frontend to enable debugging
	EnableDebuggingCommand CommandType = iota

	// DisableDebuggingCommand is a command from the frontend to disable debugging
	DisableDebuggingCommand = iota

	// ContinueCommand is a command from the front end to go to the next step
	ContinueCommand = iota
)

// CommandInterface is the interface for commands received from the user interface
type CommandInterface interface{}

// Command represents a command from the client
type Command struct {
	Type CommandType
}

// UpdateType is the type of update message for the frontend
type UpdateType uint

const (

	// RequestUpdate sends a new request to the client
	RequestUpdate UpdateType = iota

	// ResponseUpdate sends a new request to the client
	ResponseUpdate = iota

	// DebuggingToggleUpdate tells the client that debugging was turned on/off
	DebuggingToggleUpdate = iota

	// InitialUpdate tells a newly joined client everything he needs to know
	InitialUpdate = iota
)

// UpdateInterface represents an update message
type UpdateInterface interface{}

// InitialUpdateMessage is sent to the client upon initial connection
type InitialUpdateMessage struct {
	Type             UpdateType
	DebuggingEnabled bool
}

// NewInitialUpdateMessage creates a new update message
func NewInitialUpdateMessage(debuggingEnabled bool) InitialUpdateMessage {
	return InitialUpdateMessage{
		Type:             InitialUpdate,
		DebuggingEnabled: debuggingEnabled,
	}
}

// RequestUpdateMessage represents a new request update
type RequestUpdateMessage struct {
	Type UpdateType

	//TODO: decompose this
	Request    string
	RequestURI string
	Host       string
}

// NewRequestUpdateMessage creates a new update
func NewRequestUpdateMessage(request string, host string, requestURI string) RequestUpdateMessage {
	return RequestUpdateMessage{
		Type:       RequestUpdate,
		Request:    request,
		RequestURI: requestURI,
		Host:       host,
	}
}

// ResponseUpdateMessage represents a new request update
type ResponseUpdateMessage struct {
	Type UpdateType

	//TODO: decompose this
	Response string
}

// NewResponseUpdateMessage creates a new update
func NewResponseUpdateMessage(response string) ResponseUpdateMessage {
	return ResponseUpdateMessage{
		Type:     ResponseUpdate,
		Response: response,
	}
}

// DebuggingToggleMessage tells the client that debugging was turned on/off
type DebuggingToggleMessage struct {
	Type             UpdateType
	DebuggingEnabled bool
}

// NewDebuggingToggleMessage creates a new DebuggingToggleMessage
func NewDebuggingToggleMessage(debuggingEnabled bool) DebuggingToggleMessage {
	return DebuggingToggleMessage{
		Type:             DebuggingToggleUpdate,
		DebuggingEnabled: debuggingEnabled,
	}
}
