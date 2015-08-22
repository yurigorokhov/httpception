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
	NewRequestUpdate        UpdateType = iota
	NewResponseUpdate                  = iota
	DebuggingEnabledUpdate             = iota
	DebuggingDisabledUpdate            = iota
)

type Update struct {
	Type  UpdateType
	Value string
}
