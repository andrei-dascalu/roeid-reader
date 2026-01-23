package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/application"
	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/infrastructure"
)

// Romanian CEI application AID (D2 76 00 01 24 01)
var ceiAID = []byte{0xD2, 0x76, 0x00, 0x01, 0x24, 0x01}

func main() {
	fmt.Println("=== Romanian eID Reader ===")
	fmt.Println()

	// Initialize smart card service with APDU logging
	transport := infrastructure.NewPCSCTransport()
	logger := infrastructure.NewAPDULogger(os.Stdout)
	service := application.NewSmartCardService(transport, logger)

	// Connect to card (logs reader selection, ATR, protocol)
	fmt.Println("Connecting to smart card...")
	if err := service.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer service.Disconnect()

	fmt.Println()

	// SELECT Romanian eID application
	fmt.Printf("Selecting CEI application (AID: %02X)...\n", ceiAID)
	resp, err := service.SelectApplication(ceiAID)
	if err != nil {
		log.Fatalf("Failed to SELECT application: %v", err)
	}
	if len(resp.Data) > 0 {
		fmt.Printf("FCI data: %02X\n", resp.Data)
	}

	fmt.Println()

	// Prompt user for PIN1
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter PIN1: ")
	pinInput, _ := reader.ReadString('\n')
	pin := []byte(strings.TrimSpace(pinInput))

	// VERIFY PIN1 (reference 0x01)
	fmt.Println("Verifying PIN1...")
	if err := service.VerifyPIN(pin, 0x01); err != nil {
		log.Fatalf("PIN1 verification failed: %v", err)
	}

	fmt.Println()
	fmt.Println("✓ PIN1 verified successfully!")
	fmt.Println()

	// TODO: Implement PACE protocol phases to establish secure messaging
	// TODO: Read and parse protected identity data
	fmt.Println("Next steps: Implement PACE protocol (see IMPLEMENTATION_ROADMAP.md)")
}
