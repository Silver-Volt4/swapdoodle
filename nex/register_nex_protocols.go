package nex

import (
	message_delivery "github.com/PretendoNetwork/nex-protocols-go/v2/message-delivery"
	"github.com/silver-volt4/swapdoodle/globals"
	nex_message_delivery "github.com/silver-volt4/swapdoodle/nex/message-delivery"
)

func registerNEXProtocols() {
	messageDeliveryProtocol := message_delivery.NewProtocol()
	messageDeliveryProtocol.DeliverMessage = nex_message_delivery.DeliverMessage
	globals.HppServer.RegisterServiceProtocol(messageDeliveryProtocol)
}
