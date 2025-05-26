package nex_datastore_swapdoodle

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_db "github.com/silver-volt4/swapdoodle/database/datastore"
	"github.com/silver-volt4/swapdoodle/globals"
)

func GetSpecificMetaV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetSpecificMetaParamV1) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	metas, nErr := datastore_db.GetSpecificMetaByIDs(connection.PID(), param.DataIDs)
	if nErr != nil {
		globals.Logger.Error(nErr.Error())
		return nil, nErr
	}

	resStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	metas.WriteTo(resStream)

	res := nex.NewRMCSuccess(endpoint, resStream.Bytes())
	res.ProtocolID = datastore.ProtocolID
	res.MethodID = datastore.MethodGetSpecificMetaV1
	res.CallID = callID

	return res, nil
}
