package infrastructure

import (
	"fmt"
	"io"
	"time"
)

// APDULogger logs APDU commands and responses
type APDULogger struct {
	out io.Writer
}

// NewAPDULogger creates a new APDU logger
func NewAPDULogger(out io.Writer) *APDULogger {
	return &APDULogger{out: out}
}

// LogCommand logs an outgoing APDU command
func (l *APDULogger) LogCommand(data []byte) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Fprintf(l.out, "[%s] → APDU Command (%d bytes): %02X\n", timestamp, len(data), data)
}

// LogResponse logs an incoming APDU response
func (l *APDULogger) LogResponse(data []byte) {
	timestamp := time.Now().Format("15:04:05.000")
	if len(data) > 2 {
		fmt.Fprintf(l.out, "[%s] ← APDU Response (%d bytes): %02X (SW: %02X%02X)\n",
			timestamp, len(data), data, data[len(data)-2], data[len(data)-1])
	} else {
		fmt.Fprintf(l.out, "[%s] ← APDU Response (%d bytes): %02X\n",
			timestamp, len(data), data)
	}
}

// LogError logs an error
func (l *APDULogger) LogError(err error) {
	timestamp := time.Now().Format("15:04:05.000")
	fmt.Fprintf(l.out, "[%s] ✗ Error: %v\n", timestamp, err)
}
