package domain

import (
	"testing"
)

func TestAPDU_Bytes_HeaderOnly(t *testing.T) {
	apdu := &APDU{
		CLA: 0x00,
		INS: 0xA4,
		P1:  0x04,
		P2:  0x00,
	}

	got := apdu.Bytes()
	want := []byte{0x00, 0xA4, 0x04, 0x00}

	if len(got) != len(want) {
		t.Errorf("Bytes() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Bytes()[%d] = %02X, want %02X", i, got[i], want[i])
		}
	}
}

func TestAPDU_Bytes_HeaderWithLe(t *testing.T) {
	apdu := &APDU{
		CLA: 0x00,
		INS: 0xB0,
		P1:  0x00,
		P2:  0x00,
		Le:  0x10,
	}

	got := apdu.Bytes()
	want := []byte{0x00, 0xB0, 0x00, 0x00, 0x10}

	if len(got) != len(want) {
		t.Errorf("Bytes() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Bytes()[%d] = %02X, want %02X", i, got[i], want[i])
		}
	}
}

func TestAPDU_Bytes_WithData(t *testing.T) {
	// SELECT command with AID
	apdu := &APDU{
		CLA:  0x00,
		INS:  0xA4,
		P1:   0x04,
		P2:   0x00,
		Data: []byte{0xD2, 0x76, 0x00, 0x01, 0x24, 0x01},
	}

	got := apdu.Bytes()
	// Expected: CLA INS P1 P2 Lc Data...
	want := []byte{0x00, 0xA4, 0x04, 0x00, 0x06, 0xD2, 0x76, 0x00, 0x01, 0x24, 0x01}

	if len(got) != len(want) {
		t.Errorf("Bytes() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Bytes()[%d] = %02X, want %02X", i, got[i], want[i])
		}
	}
}

func TestAPDU_Bytes_WithDataAndLe(t *testing.T) {
	// Le > 0 should be appended after data
	apdu := &APDU{
		CLA:  0x00,
		INS:  0xA4,
		P1:   0x04,
		P2:   0x00,
		Data: []byte{0xD2, 0x76, 0x00, 0x01, 0x24, 0x01},
		Le:   0x10, // Request 16 bytes response
	}

	got := apdu.Bytes()
	// Expected: CLA INS P1 P2 Lc Data... Le
	want := []byte{0x00, 0xA4, 0x04, 0x00, 0x06, 0xD2, 0x76, 0x00, 0x01, 0x24, 0x01, 0x10}

	if len(got) != len(want) {
		t.Errorf("Bytes() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Bytes()[%d] = %02X, want %02X", i, got[i], want[i])
		}
	}
}

func TestAPDU_Bytes_LeZeroNotAppended(t *testing.T) {
	// Le = 0 is treated as "unspecified" and not appended (Case 3 APDU)
	apdu := &APDU{
		CLA:  0x00,
		INS:  0xA4,
		P1:   0x04,
		P2:   0x00,
		Data: []byte{0xD2, 0x76, 0x00, 0x01, 0x24, 0x01},
		Le:   0x00, // Zero = not appended
	}

	got := apdu.Bytes()
	// Expected: CLA INS P1 P2 Lc Data... (no Le)
	want := []byte{0x00, 0xA4, 0x04, 0x00, 0x06, 0xD2, 0x76, 0x00, 0x01, 0x24, 0x01}

	if len(got) != len(want) {
		t.Errorf("Bytes() length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Bytes()[%d] = %02X, want %02X", i, got[i], want[i])
		}
	}
}

func TestNewResponse_Success(t *testing.T) {
	// Response with data and success status
	data := []byte{0x6F, 0x10, 0x84, 0x06, 0xD2, 0x76, 0x00, 0x01, 0x24, 0x01, 0x90, 0x00}
	resp := NewResponse(data)

	if resp.SW1 != 0x90 {
		t.Errorf("SW1 = %02X, want 0x90", resp.SW1)
	}
	if resp.SW2 != 0x00 {
		t.Errorf("SW2 = %02X, want 0x00", resp.SW2)
	}
	if !resp.IsSuccess() {
		t.Error("IsSuccess() = false, want true")
	}
	if resp.StatusCode() != 0x9000 {
		t.Errorf("StatusCode() = %04X, want 0x9000", resp.StatusCode())
	}
	if len(resp.Data) != 10 {
		t.Errorf("Data length = %d, want 10", len(resp.Data))
	}
}

func TestNewResponse_Error(t *testing.T) {
	// Security status not satisfied
	data := []byte{0x69, 0x82}
	resp := NewResponse(data)

	if resp.SW1 != 0x69 {
		t.Errorf("SW1 = %02X, want 0x69", resp.SW1)
	}
	if resp.SW2 != 0x82 {
		t.Errorf("SW2 = %02X, want 0x82", resp.SW2)
	}
	if resp.IsSuccess() {
		t.Error("IsSuccess() = true, want false")
	}
	if resp.StatusCode() != StatusSecurityAuthFailed {
		t.Errorf("StatusCode() = %04X, want %04X", resp.StatusCode(), StatusSecurityAuthFailed)
	}
	if len(resp.Data) != 0 {
		t.Errorf("Data length = %d, want 0", len(resp.Data))
	}
}

func TestNewResponse_TooShort(t *testing.T) {
	// Malformed response with only 1 byte
	data := []byte{0x90}
	resp := NewResponse(data)

	// Should return unknown error status
	if resp.SW1 != 0x6F {
		t.Errorf("SW1 = %02X, want 0x6F (unknown error)", resp.SW1)
	}
}

func TestResponse_Error(t *testing.T) {
	resp := &Response{SW1: 0x69, SW2: 0x82}
	errStr := resp.Error()

	if errStr == "" {
		t.Error("Error() returned empty string")
	}
}

func TestNewResponse_PINRetryCounter(t *testing.T) {
	// PIN retry counter: 3 retries left (0x63C3)
	data := []byte{0x63, 0xC3}
	resp := NewResponse(data)

	if resp.StatusCode() != 0x63C3 {
		t.Errorf("StatusCode() = %04X, want 0x63C3", resp.StatusCode())
	}
	if resp.IsSuccess() {
		t.Error("IsSuccess() = true, want false for retry counter")
	}
}
