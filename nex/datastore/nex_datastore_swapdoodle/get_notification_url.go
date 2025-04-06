package nex_datastore_swapdoodle

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/silver-volt4/swapdoodle/globals"
)

func GetNotificationURL(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetNotificationURLParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.HppServer.LibraryVersions(), globals.HppServer.ByteStreamSettings())

	urlInfo := datastore_types.NewDataStoreReqGetNotificationURLInfo()

	urlInfo.URL = types.NewString("https://example.com")
	urlInfo.Key = types.NewString("some/key")
	urlInfo.Query = types.NewString("?test=test")
	urlInfo.RootCACert = types.NewBuffer(nil)

	urlInfo.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.HppServer, rmcResponseStream.Bytes()) // rmcResponseStream.Bytes()
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodGetNotificationURL
	rmcResponse.CallID = callID

	// I'm actually returning an error here because otherwise it would
	// connect to example.com every second request and that would slow down development
	return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "asdf")

	//return rmcResponse, nil
}
