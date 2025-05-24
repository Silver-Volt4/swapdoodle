package nex_datastore_swapdoodle

import (
	"fmt"
	"time"

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

	bucket := globals.DatastoreCommon.S3Bucket
	key := fmt.Sprintf("%s/%d", "notifications", packet.Sender().PID())

	url, err := globals.DatastoreCommon.S3Presigner.GetObject(bucket, key, time.Hour*24*7)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "change_error")
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.HppServer.LibraryVersions(), globals.HppServer.ByteStreamSettings())

	urlInfo := datastore_types.NewDataStoreReqGetNotificationURLInfo()

	urlInfo.URL = types.NewString(url.Host)
	urlInfo.Key = types.NewString(url.Path)
	urlInfo.Query = types.NewString("?" + url.Query().Encode())
	urlInfo.RootCACert = types.NewBuffer(nil)

	urlInfo.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.HppServer, rmcResponseStream.Bytes()) // rmcResponseStream.Bytes()
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodGetNotificationURL
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
