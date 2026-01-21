package domain

import "fmt"

// APDU represents a smart card command/response (ISO/IEC 7816-4)
type APDU struct {
	CLA  byte   // Class
	INS  byte   // Instruction
	P1   byte   // Parameter 1
	P2   byte   // Parameter 2
	Data []byte // Command data
	Le   byte   // Expected response length (0 = unspecified)
}

// Bytes serializes APDU to ISO/IEC 7816-4 format
func (a *APDU) Bytes() []byte {
	cmd := []byte{a.CLA, a.INS, a.P1, a.P2}
	if len(a.Data) == 0 {
		// No data, just header + Le
		if a.Le > 0 {
			cmd = append(cmd, a.Le)
		}
		return cmd
	}
	// With data
	cmd = append(cmd, byte(len(a.Data)))
	cmd = append(cmd, a.Data...)
	if a.Le > 0 {
		cmd = append(cmd, a.Le)
	}
	return cmd
}

// Response represents a smart card response
type Response struct {
	Data []byte // Response data
	SW1  byte   // Status word 1
	SW2  byte   // Status word 2
}

// StatusCode returns the 16-bit status word (SW1 << 8 | SW2)
func (r *Response) StatusCode() uint16 {
	return (uint16(r.SW1) << 8) | uint16(r.SW2)
}

// IsSuccess checks if status is 0x9000 (success)
func (r *Response) IsSuccess() bool {
	return r.StatusCode() == 0x9000
}

// Error implements the error interface
func (r *Response) Error() string {
	return fmt.Sprintf("APDU error: SW1=%02X SW2=%02X", r.SW1, r.SW2)
}

// NewResponse creates a Response from raw bytes (last 2 bytes are status words)
func NewResponse(data []byte) *Response {
	if len(data) < 2 {
		return &Response{Data: data, SW1: 0x6F, SW2: 0x00} // Unknown error
	}
	return &Response{
		Data: data[:len(data)-2],
		SW1:  data[len(data)-2],
		SW2:  data[len(data)-1],
	}
}
