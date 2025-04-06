package nex

import (
	"fmt"
	"os"

	"github.com/PretendoNetwork/nex-go"
	"github.com/silver-volt4/swapdoodle/globals"
)

const ACCESS_KEY = "76f26496" // The glorious result of ~7 hours of brute forcing!

func StartHppServer() {
	globals.HppServer = nex.NewServer()
	globals.HppServer.SetDefaultNEXVersion(nex.NewPatchedNEXVersion(3, 8, 3, "AMAJ"))
	globals.HppServer.SetKerberosPassword(globals.KerberosPassword)
	globals.HppServer.SetAccessKey(ACCESS_KEY)

	globals.HppServer.On("Data", func(packet *nex.HPPPacket) {
		request := packet.RMCRequest()

		fmt.Println("== Swapdoodle - HPP ==")
		fmt.Printf("Protocol ID: %d\n", request.ProtocolID())
		fmt.Printf("Method ID: %d\n", request.MethodID())
		fmt.Println("===============")
	})

	// * Register the common handlers first so that they can be overridden if needed
	registerCommonProtocols()
	registerNEXProtocols()

	globals.HppServer.HPPListen(fmt.Sprintf(":%s", os.Getenv("PN_SD_HPP_SERVER_PORT")))
}
