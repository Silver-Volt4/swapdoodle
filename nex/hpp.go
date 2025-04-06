package nex

import (
	"fmt"
	"os"
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/silver-volt4/swapdoodle/globals"
)

const ACCESS_KEY = "76f26496" // The glorious result of ~7 hours of brute forcing!

func StartHppServer() {
	globals.HppServer = nex.NewHPPServer()

	globals.HppServer.AccountDetailsByPID = globals.AccountDetailsByPID
	globals.HppServer.AccountDetailsByUsername = globals.AccountDetailsByUsername
	globals.HppServer.LibraryVersions().SetDefault(nex.NewLibraryVersion(3, 8, 3))
	globals.HppServer.SetAccessKey(ACCESS_KEY)

	globals.HppServer.OnData(func(packet nex.PacketInterface) {
		request := packet.RMCMessage()

		fmt.Println("=== Swapdoodle - HPP ===")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID)
		fmt.Printf("Method ID: %d\n", request.MethodID)
		fmt.Println("==================")
	})

	// * Register the common handlers first so that they can be overridden if needed
	registerProtocols()

	port, _ := strconv.Atoi(os.Getenv("PN_SD_HPP_SERVER_PORT"))

	globals.HppServer.Listen(port)
}
