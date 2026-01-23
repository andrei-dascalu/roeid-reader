package domain

// Status code constants (ISO/IEC 7816-4)
const (
	// Successful execution
	StatusSuccess uint16 = 0x9000

	// Execution errors (0x61xx, 0x62xx, 0x63xx)
	StatusMoreData         uint16 = 0x6100 // More data available
	StatusWarning          uint16 = 0x6200 // Warning
	StatusWarningCorrupted uint16 = 0x6281 // File corrupted
	StatusWarningEOF       uint16 = 0x6282 // End of file reached

	// Execution errors (0x64xx, 0x65xx, 0x66xx)
	StatusExecutionError  uint16 = 0x6400 // Execution error
	StatusPersistentError uint16 = 0x6500 // Persistent error
	StatusSecurityError   uint16 = 0x6600 // Security error

	// Client errors (0x67xx, 0x68xx, 0x69xx, 0x6Axx)
	StatusLengthError       uint16 = 0x6700 // Wrong length
	StatusFunctionNotFound  uint16 = 0x6881 // Logical channel not supported
	StatusLogicalChannelErr uint16 = 0x6882 // Secure messaging not supported
	StatusKeyReferenceErr   uint16 = 0x6A86 // Incorrect parameters (P1/P2)
	StatusFileNotFound      uint16 = 0x6A82 // File or application not found
	StatusInstructionErr    uint16 = 0x6D00 // Instruction code not supported
	StatusCLAErr            uint16 = 0x6E00 // Class not supported

	// Authentication errors
	StatusSecurityAuthFailed uint16 = 0x6982 // Security status not satisfied
	StatusIncorrectPIN       uint16 = 0x6983 // Authentication method blocked
	StatusIncorrectKey       uint16 = 0x6984 // Reference key in use

	// PIN retry counter (0x63Cx where x = remaining tries)
	StatusPINRetryMask uint16 = 0x63C0 // Mask for PIN retry counter
)

// TransportErrorCode represents PC/SC transport-level errors
type TransportErrorCode int

const (
	ErrNoContext TransportErrorCode = iota + 1
	ErrNoReaders
	ErrNoCard
	ErrCardRemoved
	ErrConnectionLost
	ErrTransmissionFailed
	ErrProtocolMismatch
	ErrReaderBusy
	ErrTimeout
)

// TransportError represents a PC/SC transport-level error (not APDU status)
type TransportError struct {
	Code    TransportErrorCode
	Message string
	Cause   error
}

// Error implements the error interface
func (e *TransportError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *TransportError) Unwrap() error {
	return e.Cause
}

// NewTransportError creates a new transport error
func NewTransportError(code TransportErrorCode, message string, cause error) *TransportError {
	return &TransportError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// StatusError wraps a status code with context
type StatusError struct {
	Code   uint16
	SW1    byte
	SW2    byte
	Detail string
}

// Error implements the error interface
func (e *StatusError) Error() string {
	if e.Detail != "" {
		return e.Detail
	}
	return e.String()
}

// String returns a human-readable status code
func (e *StatusError) String() string {
	switch e.Code {
	case StatusSuccess:
		return "Success"
	case StatusMoreData:
		return "More data available"
	case StatusWarning:
		return "Warning"
	case StatusWarningCorrupted:
		return "File corrupted"
	case StatusWarningEOF:
		return "End of file reached"
	case StatusExecutionError:
		return "Execution error"
	case StatusPersistentError:
		return "Persistent memory error"
	case StatusSecurityError:
		return "Security error"
	case StatusLengthError:
		return "Wrong length"
	case StatusSecurityAuthFailed:
		return "Security status not satisfied (incorrect PIN/CAN?)"
	case StatusIncorrectPIN:
		return "Authentication method blocked"
	case StatusInstructionErr:
		return "Instruction code not supported"
	case StatusCLAErr:
		return "Class not supported"
	default:
		return "Unknown error"
	}
}

// NewStatusError creates a StatusError from a Response
func NewStatusError(resp *Response) *StatusError {
	return &StatusError{
		Code: resp.StatusCode(),
		SW1:  resp.SW1,
		SW2:  resp.SW2,
	}
}
