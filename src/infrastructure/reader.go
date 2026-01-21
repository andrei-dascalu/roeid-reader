package infrastructure

import (
	"fmt"
	"log"
	"strings"

	"github.com/ebfe/scard"
)

const LeaveCard scard.Disposition = scard.LeaveCard

func Connect() *scard.Context {
	ctx, err := scard.EstablishContext()
	if err != nil {
		log.Fatal("Failed to establish context:", err)
	}

	return ctx
}

func listReaders(ctx *scard.Context) []string {
	readers, err := ctx.ListReaders()
	if err != nil {
		log.Fatal("Failed to list readers:", err)
	}

	return readers
}

func SelectReader(ctx *scard.Context) string {
	readers := listReaders(ctx)

	if len(readers) == 0 {
		log.Fatal("No smart card readers found")
	}

	selected := readers[0]

	for _, r := range readers {
		fmt.Println("Reader: ", r)

		if strings.Contains(r, "Generic EMV") {
			return r
		}
	}

	return selected
}

func ConnectCard(ctx *scard.Context, readerName string) *scard.Card {
	card, err := ctx.Connect(readerName, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		log.Fatal("Failed to connect to card:", err)
	}

	return card
}

func DiagnoseCard(card *scard.Card) {
	status, err := card.Status()
	if err != nil {
		log.Fatal("Failed to get card status:", err)
	}

	fmt.Println("Card ATR:", fmt.Sprintf("% X", status.Atr))
	fmt.Println("Active Protocol:", status.ActiveProtocol)
}
