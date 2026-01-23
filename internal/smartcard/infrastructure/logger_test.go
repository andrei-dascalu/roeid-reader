package infrastructure

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestNewAPDULogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	if logger == nil {
		t.Fatal("NewAPDULogger returned nil")
	}
	if !logger.enabled {
		t.Error("Logger should be enabled by default")
	}
}

func TestAPDULogger_SetEnabled(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	logger.SetEnabled(false)
	logger.LogInfo("test")

	if buf.Len() != 0 {
		t.Error("Disabled logger should not write output")
	}

	logger.SetEnabled(true)
	logger.LogInfo("test")

	if buf.Len() == 0 {
		t.Error("Enabled logger should write output")
	}
}

func TestAPDULogger_LogConnect(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	logger.LogConnect("Test Reader", "T=1")

	output := buf.String()
	if !strings.Contains(output, "●") {
		t.Error("LogConnect should contain connection symbol")
	}
	if !strings.Contains(output, "Test Reader") {
		t.Error("LogConnect should contain reader name")
	}
	if !strings.Contains(output, "T=1") {
		t.Error("LogConnect should contain protocol")
	}
}

func TestAPDULogger_LogDisconnect(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	logger.LogDisconnect()

	output := buf.String()
	if !strings.Contains(output, "○") {
		t.Error("LogDisconnect should contain disconnection symbol")
	}
	if !strings.Contains(output, "Disconnected") {
		t.Error("LogDisconnect should contain 'Disconnected'")
	}
}

func TestAPDULogger_LogReaderSelected(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	readers := []string{"Reader 1", "Reader 2", "Reader 3"}
	logger.LogReaderSelected("Reader 2", readers)

	output := buf.String()
	if !strings.Contains(output, "◆") {
		t.Error("LogReaderSelected should contain selection symbol")
	}
	if !strings.Contains(output, "Reader 2") {
		t.Error("LogReaderSelected should contain selected reader")
	}
	if !strings.Contains(output, "3 available") {
		t.Error("LogReaderSelected should show count of available readers")
	}
}

func TestAPDULogger_LogATR(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	atr := []byte{0x3B, 0x8C, 0x80, 0x01}
	logger.LogATR(atr)

	output := buf.String()
	if !strings.Contains(output, "◇") {
		t.Error("LogATR should contain ATR symbol")
	}
	if !strings.Contains(output, "3B8C8001") {
		t.Error("LogATR should contain hex-encoded ATR")
	}
}

func TestAPDULogger_LogCommand_SELECT(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	// SELECT command
	cmd := []byte{0x00, 0xA4, 0x04, 0x00, 0x06, 0xD2, 0x76, 0x00, 0x01, 0x24, 0x01}
	logger.LogCommand(cmd)

	output := buf.String()
	if !strings.Contains(output, "→") {
		t.Error("LogCommand should contain arrow symbol")
	}
	if !strings.Contains(output, "SELECT") {
		t.Error("LogCommand should identify SELECT instruction")
	}
	if !strings.Contains(output, "CLA=00") {
		t.Error("LogCommand should show CLA byte")
	}
	if !strings.Contains(output, "INS=A4") {
		t.Error("LogCommand should show INS byte")
	}
}

func TestAPDULogger_LogCommand_VERIFY(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	// VERIFY command
	cmd := []byte{0x00, 0x20, 0x00, 0x01, 0x04, 0x31, 0x32, 0x33, 0x34}
	logger.LogCommand(cmd)

	output := buf.String()
	if !strings.Contains(output, "VERIFY") {
		t.Error("LogCommand should identify VERIFY instruction")
	}
}

func TestAPDULogger_LogCommand_UnknownINS(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	// Unknown instruction
	cmd := []byte{0x00, 0xFF, 0x00, 0x00}
	logger.LogCommand(cmd)

	output := buf.String()
	if !strings.Contains(output, "INS_FF") {
		t.Error("LogCommand should show INS_XX for unknown instructions")
	}
}

func TestAPDULogger_LogCommand_ShortAPDU(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	// Less than 4 bytes
	cmd := []byte{0x00, 0xA4}
	logger.LogCommand(cmd)

	output := buf.String()
	if !strings.Contains(output, "APDU Command") {
		t.Error("Short APDU should use generic format")
	}
}

func TestAPDULogger_LogResponse_Success(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	// Success response with data
	resp := []byte{0x6F, 0x10, 0x90, 0x00}
	logger.LogResponse(resp)

	output := buf.String()
	if !strings.Contains(output, "←") {
		t.Error("LogResponse should contain arrow symbol")
	}
	if !strings.Contains(output, "SW=9000") {
		t.Error("LogResponse should show status word")
	}
	if !strings.Contains(output, "OK") {
		t.Error("LogResponse should show 'OK' for success")
	}
	if !strings.Contains(output, "2 data bytes") {
		t.Error("LogResponse should show data byte count")
	}
}

func TestAPDULogger_LogResponse_AuthFailed(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	resp := []byte{0x69, 0x82}
	logger.LogResponse(resp)

	output := buf.String()
	if !strings.Contains(output, "SW=6982") {
		t.Error("LogResponse should show status word")
	}
	if !strings.Contains(output, "Auth failed") {
		t.Error("LogResponse should show 'Auth failed' for 6982")
	}
}

func TestAPDULogger_LogResponse_PINRetry(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	// PIN retry counter: 2 retries left
	resp := []byte{0x63, 0xC2}
	logger.LogResponse(resp)

	output := buf.String()
	if !strings.Contains(output, "SW=63C2") {
		t.Error("LogResponse should show status word")
	}
	if !strings.Contains(output, "2 retries left") {
		t.Error("LogResponse should show retry count")
	}
}

func TestAPDULogger_LogResponse_ShortResponse(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	// Less than 2 bytes
	resp := []byte{0x90}
	logger.LogResponse(resp)

	output := buf.String()
	if !strings.Contains(output, "1 bytes") {
		t.Error("Short response should use generic format")
	}
}

func TestAPDULogger_LogError(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	err := errors.New("test error message")
	logger.LogError(err)

	output := buf.String()
	if !strings.Contains(output, "✗") {
		t.Error("LogError should contain error symbol")
	}
	if !strings.Contains(output, "test error message") {
		t.Error("LogError should contain error message")
	}
}

func TestAPDULogger_LogInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	logger.LogInfo("Test message with %s", "formatting")

	output := buf.String()
	if !strings.Contains(output, "ℹ") {
		t.Error("LogInfo should contain info symbol")
	}
	if !strings.Contains(output, "Test message with formatting") {
		t.Error("LogInfo should support format strings")
	}
}

func TestAPDULogger_Timestamp(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)

	logger.LogInfo("test")

	output := buf.String()
	// Timestamp format: HH:MM:SS.mmm
	if !strings.Contains(output, ":") {
		t.Error("Output should contain timestamp with colons")
	}
	if !strings.Contains(output, ".") {
		t.Error("Output should contain timestamp with milliseconds")
	}
}

func TestAPDULogger_InstructionName(t *testing.T) {
	logger := NewAPDULogger(&bytes.Buffer{})

	tests := []struct {
		ins  byte
		name string
	}{
		{0xA4, "SELECT"},
		{0x20, "VERIFY"},
		{0xB0, "READ BINARY"},
		{0xB2, "READ RECORD"},
		{0xCA, "GET DATA"},
		{0xD6, "UPDATE BINARY"},
		{0x82, "EXTERNAL AUTHENTICATE"},
		{0x84, "GET CHALLENGE"},
		{0x86, "GENERAL AUTHENTICATE"},
		{0x88, "INTERNAL AUTHENTICATE"},
		{0x22, "MANAGE SECURITY ENVIRONMENT"},
		{0x2A, "PERFORM SECURITY OPERATION"},
		{0xFF, "INS_FF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := logger.instructionName(tt.ins)
			if got != tt.name {
				t.Errorf("instructionName(0x%02X) = %q, want %q", tt.ins, got, tt.name)
			}
		})
	}
}

func TestAPDULogger_Disabled(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewAPDULogger(buf)
	logger.SetEnabled(false)

	// All log methods should be no-ops when disabled
	logger.LogConnect("reader", "T=1")
	logger.LogDisconnect()
	logger.LogReaderSelected("reader", []string{"reader"})
	logger.LogATR([]byte{0x3B})
	logger.LogCommand([]byte{0x00, 0xA4, 0x04, 0x00})
	logger.LogResponse([]byte{0x90, 0x00})
	logger.LogError(errors.New("test"))
	logger.LogInfo("test")

	if buf.Len() != 0 {
		t.Errorf("Disabled logger should not produce output, got %d bytes", buf.Len())
	}
}
