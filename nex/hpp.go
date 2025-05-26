package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/silver-volt4/swapdoodle/globals"
)

func StartHppServer() {
	globals.HppServer = nex.NewHPPServer()

	globals.HppServer.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.HppServer.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.HppServer.LibraryVersions().SetDefault(globals.LibraryVersion)
	globals.HppServer.SetAccessKey(globals.HPP_ACCESS_KEY)

	globals.HppServer.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== Swapdoodle - HPP ===")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("==================")
	})

	registerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_SD_HPP_SERVER_PORT"))

	globals.HppServer.Listen(port)
}
