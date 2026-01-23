package infrastructure

import (
	"fmt"
	"io"
	"time"

	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/domain"
)

// APDULogger logs APDU commands, responses, and connection lifecycle events
type APDULogger struct {
	out     io.Writer
	enabled bool
}

// NewAPDULogger creates a new APDU logger
func NewAPDULogger(out io.Writer) *APDULogger {
	return &APDULogger{out: out, enabled: true}
}

// SetEnabled enables or disables logging
func (l *APDULogger) SetEnabled(enabled bool) {
	l.enabled = enabled
}

func (l *APDULogger) timestamp() string {
	return time.Now().Format("15:04:05.000")
}

// LogConnect logs a successful card connection
func (l *APDULogger) LogConnect(reader string, protocol string) {
	if !l.enabled {
		return
	}
	fmt.Fprintf(l.out, "[%s] ● Connected to reader: %s (protocol: %s)\n",
		l.timestamp(), reader, protocol)
}

// LogDisconnect logs a card disconnection
func (l *APDULogger) LogDisconnect() {
	if !l.enabled {
		return
	}
	fmt.Fprintf(l.out, "[%s] ○ Disconnected from card\n", l.timestamp())
}

// LogReaderSelected logs reader selection
func (l *APDULogger) LogReaderSelected(reader string, available []string) {
	if !l.enabled {
		return
	}
	fmt.Fprintf(l.out, "[%s] ◆ Selected reader: %s (from %d available)\n",
		l.timestamp(), reader, len(available))
}

// LogATR logs the Answer To Reset
func (l *APDULogger) LogATR(atr []byte) {
	if !l.enabled {
		return
	}
	fmt.Fprintf(l.out, "[%s] ◇ Card ATR: %02X\n", l.timestamp(), atr)
}

// LogCommand logs an outgoing APDU command with parsed header
func (l *APDULogger) LogCommand(data []byte) {
	if !l.enabled {
		return
	}
	if len(data) >= 4 {
		// Parse APDU header for readability
		cla, ins, p1, p2 := data[0], data[1], data[2], data[3]
		insName := l.instructionName(ins)
		fmt.Fprintf(l.out, "[%s] → APDU %s (CLA=%02X INS=%02X P1=%02X P2=%02X, %d bytes): %02X\n",
			l.timestamp(), insName, cla, ins, p1, p2, len(data), data)
	} else {
		fmt.Fprintf(l.out, "[%s] → APDU Command (%d bytes): %02X\n",
			l.timestamp(), len(data), data)
	}
}

// LogResponse logs an incoming APDU response with status interpretation
func (l *APDULogger) LogResponse(data []byte) {
	if !l.enabled {
		return
	}
	if len(data) >= 2 {
		sw1, sw2 := data[len(data)-2], data[len(data)-1]
		statusCode := (uint16(sw1) << 8) | uint16(sw2)
		statusDesc := l.statusDescription(statusCode)
		dataLen := len(data) - 2
		fmt.Fprintf(l.out, "[%s] ← APDU Response (%d data bytes, SW=%04X %s): %02X\n",
			l.timestamp(), dataLen, statusCode, statusDesc, data)
	} else {
		fmt.Fprintf(l.out, "[%s] ← APDU Response (%d bytes): %02X\n",
			l.timestamp(), len(data), data)
	}
}

// LogError logs an error
func (l *APDULogger) LogError(err error) {
	if !l.enabled {
		return
	}
	fmt.Fprintf(l.out, "[%s] ✗ Error: %v\n", l.timestamp(), err)
}

// LogInfo logs an informational message
func (l *APDULogger) LogInfo(format string, args ...any) {
	if !l.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.out, "[%s] ℹ %s\n", l.timestamp(), msg)
}

// instructionName returns a human-readable name for common INS bytes (ISO/IEC 7816-4)
func (l *APDULogger) instructionName(ins byte) string {
	switch ins {
	case 0xA4:
		return "SELECT"
	case 0x20:
		return "VERIFY"
	case 0xB0:
		return "READ BINARY"
	case 0xB2:
		return "READ RECORD"
	case 0xCA:
		return "GET DATA"
	case 0xD6:
		return "UPDATE BINARY"
	case 0x82:
		return "EXTERNAL AUTHENTICATE"
	case 0x84:
		return "GET CHALLENGE"
	case 0x86:
		return "GENERAL AUTHENTICATE"
	case 0x88:
		return "INTERNAL AUTHENTICATE"
	case 0x22:
		return "MANAGE SECURITY ENVIRONMENT"
	case 0x2A:
		return "PERFORM SECURITY OPERATION"
	default:
		return fmt.Sprintf("INS_%02X", ins)
	}
}

// statusDescription returns a human-readable status description
func (l *APDULogger) statusDescription(status uint16) string {
	switch status {
	case domain.StatusSuccess:
		return "OK"
	case domain.StatusMoreData:
		return "More data"
	case domain.StatusWarningEOF:
		return "EOF"
	case domain.StatusSecurityAuthFailed:
		return "Auth failed"
	case domain.StatusIncorrectPIN:
		return "PIN blocked"
	case domain.StatusFileNotFound:
		return "Not found"
	case domain.StatusLengthError:
		return "Wrong length"
	case domain.StatusInstructionErr:
		return "INS not supported"
	case domain.StatusCLAErr:
		return "CLA not supported"
	default:
		// Check for PIN retry counter (0x63Cx)
		if status&0xFFF0 == domain.StatusPINRetryMask {
			retries := status & 0x000F
			return fmt.Sprintf("%d retries left", retries)
		}
		return ""
	}
}
