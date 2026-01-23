package domain

import (
	"errors"
	"testing"
)

func TestStatusError_Error(t *testing.T) {
	tests := []struct {
		name   string
		err    *StatusError
		expect string
	}{
		{
			name:   "success",
			err:    &StatusError{Code: StatusSuccess, SW1: 0x90, SW2: 0x00},
			expect: "Success",
		},
		{
			name:   "security auth failed",
			err:    &StatusError{Code: StatusSecurityAuthFailed, SW1: 0x69, SW2: 0x82},
			expect: "Security status not satisfied (incorrect PIN/CAN?)",
		},
		{
			name:   "with detail override",
			err:    &StatusError{Code: StatusSecurityAuthFailed, SW1: 0x69, SW2: 0x82, Detail: "Custom error"},
			expect: "Custom error",
		},
		{
			name:   "file not found",
			err:    &StatusError{Code: StatusFileNotFound, SW1: 0x6A, SW2: 0x82},
			expect: "Unknown error", // Not mapped in String() yet
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expect {
				t.Errorf("Error() = %q, want %q", got, tt.expect)
			}
		})
	}
}

func TestStatusError_String(t *testing.T) {
	tests := []struct {
		code   uint16
		expect string
	}{
		{StatusSuccess, "Success"},
		{StatusMoreData, "More data available"},
		{StatusWarningEOF, "End of file reached"},
		{StatusSecurityAuthFailed, "Security status not satisfied (incorrect PIN/CAN?)"},
		{StatusIncorrectPIN, "Authentication method blocked"},
		{StatusLengthError, "Wrong length"},
		{StatusInstructionErr, "Instruction code not supported"},
		{StatusCLAErr, "Class not supported"},
		{0xFFFF, "Unknown error"},
	}

	for _, tt := range tests {
		t.Run(tt.expect, func(t *testing.T) {
			err := &StatusError{Code: tt.code}
			got := err.String()
			if got != tt.expect {
				t.Errorf("String() = %q, want %q", got, tt.expect)
			}
		})
	}
}

func TestNewStatusError(t *testing.T) {
	resp := &Response{
		Data: []byte{0x01, 0x02},
		SW1:  0x69,
		SW2:  0x82,
	}

	err := NewStatusError(resp)

	if err.Code != StatusSecurityAuthFailed {
		t.Errorf("Code = %04X, want %04X", err.Code, StatusSecurityAuthFailed)
	}
	if err.SW1 != 0x69 {
		t.Errorf("SW1 = %02X, want 0x69", err.SW1)
	}
	if err.SW2 != 0x82 {
		t.Errorf("SW2 = %02X, want 0x82", err.SW2)
	}
}

func TestTransportError_Error(t *testing.T) {
	tests := []struct {
		name   string
		err    *TransportError
		expect string
	}{
		{
			name:   "without cause",
			err:    &TransportError{Code: ErrNoCard, Message: "no card inserted"},
			expect: "no card inserted",
		},
		{
			name:   "with cause",
			err:    &TransportError{Code: ErrNoCard, Message: "no card inserted", Cause: errors.New("PC/SC error")},
			expect: "no card inserted: PC/SC error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expect {
				t.Errorf("Error() = %q, want %q", got, tt.expect)
			}
		})
	}
}

func TestTransportError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &TransportError{Code: ErrNoReaders, Message: "test", Cause: cause}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestTransportError_Unwrap_NoCause(t *testing.T) {
	err := &TransportError{Code: ErrNoReaders, Message: "test"}

	unwrapped := err.Unwrap()
	if unwrapped != nil {
		t.Errorf("Unwrap() = %v, want nil", unwrapped)
	}
}

func TestNewTransportError(t *testing.T) {
	cause := errors.New("scard error")
	err := NewTransportError(ErrTransmissionFailed, "transmission failed", cause)

	if err.Code != ErrTransmissionFailed {
		t.Errorf("Code = %v, want %v", err.Code, ErrTransmissionFailed)
	}
	if err.Message != "transmission failed" {
		t.Errorf("Message = %q, want %q", err.Message, "transmission failed")
	}
	if err.Cause != cause {
		t.Errorf("Cause = %v, want %v", err.Cause, cause)
	}
}

func TestTransportErrorCode_Values(t *testing.T) {
	// Verify error codes are unique and non-zero
	codes := []TransportErrorCode{
		ErrNoContext,
		ErrNoReaders,
		ErrNoCard,
		ErrCardRemoved,
		ErrConnectionLost,
		ErrTransmissionFailed,
		ErrProtocolMismatch,
		ErrReaderBusy,
		ErrTimeout,
	}

	seen := make(map[TransportErrorCode]bool)
	for _, code := range codes {
		if code == 0 {
			t.Errorf("TransportErrorCode should not be zero")
		}
		if seen[code] {
			t.Errorf("Duplicate TransportErrorCode: %d", code)
		}
		seen[code] = true
	}
}

func TestStatusConstants(t *testing.T) {
	// Verify ISO/IEC 7816-4 status codes are correctly defined
	tests := []struct {
		name   string
		code   uint16
		expect uint16
	}{
		{"Success", StatusSuccess, 0x9000},
		{"More data", StatusMoreData, 0x6100},
		{"Security auth failed", StatusSecurityAuthFailed, 0x6982},
		{"Incorrect PIN", StatusIncorrectPIN, 0x6983},
		{"File not found", StatusFileNotFound, 0x6A82},
		{"Wrong length", StatusLengthError, 0x6700},
		{"INS not supported", StatusInstructionErr, 0x6D00},
		{"CLA not supported", StatusCLAErr, 0x6E00},
		{"PIN retry mask", StatusPINRetryMask, 0x63C0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expect {
				t.Errorf("%s = %04X, want %04X", tt.name, tt.code, tt.expect)
			}
		})
	}
}
