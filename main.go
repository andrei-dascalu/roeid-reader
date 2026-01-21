package main

import (
	"fmt"

	"github.com/andrei-dascalu/roeid-reader/src/infrastructure"
)

func main() {
	cardContext := infrastructure.Connect()
	defer cardContext.Release()

	selectedReader := infrastructure.SelectReader(cardContext)

	fmt.Println("Selected: ", selectedReader)

	card := infrastructure.ConnectCard(cardContext, selectedReader)

	defer card.Disconnect(infrastructure.LeaveCard)

	infrastructure.DiagnoseCard(card)
}
