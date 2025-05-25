package nex_datastore_swapdoodle

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_db "github.com/silver-volt4/swapdoodle/database/datastore"
	"github.com/silver-volt4/swapdoodle/globals"
)

func GetSpecificMetaV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetSpecificMetaParamV1) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		fmt.Println("this failed")
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	metas, nErr := datastore_db.GetSpecificMetaByIDs(connection.PID(), param.DataIDs)

	if nErr != nil {
		fmt.Println(nErr)
		return nil, nErr
	}

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	metas.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodGetSpecificMetaV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
