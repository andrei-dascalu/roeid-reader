package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/application"
	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/infrastructure"
)

func main() {
	// Initialize smart card service
	transport := infrastructure.NewPCSCTransport()
	logger := infrastructure.NewAPDULogger(os.Stdout)
	service := application.NewSmartCardService(transport, logger)

	// Connect to card
	if err := service.Connect(); err != nil {
		log.Fatalf("Failed to connect to smart card: %v", err)
	}
	defer service.Disconnect()

	// Diagnose card
	if err := service.DiagnoseCard(); err != nil {
		log.Fatalf("Failed to diagnose card: %v", err)
	}

	fmt.Println("Card connected! Waiting for PIN1...")

	// Prompt user for PIN1
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter PIN1: ")
	pinInput, _ := reader.ReadString('\n')

	// Normalize PIN input
	pin := []byte(pinInput)
	if len(pin) > 8 { // truncate to 8 bytes
		pin = pin[:8]
	}

	// SELECT Romanian eID application
	if err := service.SelectApplication([]byte{0xD2, 0x76, 0x00, 0x01, 0x24, 0x01}); err != nil {
		log.Fatalf("Failed to SELECT application: %v", err)
	}

	// VERIFY PIN1
	if err := service.VerifyPIN(pin); err != nil {
		log.Fatalf("PIN1 verification failed: %v", err)
	}

	fmt.Println("PIN1 verified! You can now send APDUs to read personal data.")

	// TODO: Implement PACE protocol phases to establish secure messaging
	// TODO: Read and parse protected identity data
	fmt.Println("Template complete. Replace with PACE protocol implementation.")
}
