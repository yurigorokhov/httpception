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

// Command is the command message
type Command struct {
	Type  CommandType
	Value string
}

// UpdateType is the type of update message for the frontend
type UpdateType uint

const (

	// NewRequestUpdate sends a new request to the client
	NewRequestUpdate UpdateType = iota

	// NewResponseUpdate sends a new request to the client
	NewResponseUpdate = iota

	// DebuggingEnabledUpdate tells the client that debugging was turned on
	DebuggingEnabledUpdate = iota

	// DebuggingDisabledUpdate tells the client that debugging was turned off
	DebuggingDisabledUpdate = iota
)

// Update represents an update message to the client
type Update struct {
	Type  UpdateType
	Value string
}
