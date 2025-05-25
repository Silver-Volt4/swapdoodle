package nex_datastore_swapdoodle

import (
	"fmt"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_db "github.com/silver-volt4/swapdoodle/database/datastore"
	"github.com/silver-volt4/swapdoodle/globals"
)

func PrepareGetObjectV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStorePrepareGetParamV1) (*nex.RMCMessage, *nex.Error) {
	if globals.DatastoreCommon.S3Presigner == nil {
		globals.Logger.Warning("S3Presigner not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	// * Only allow the owner or recipient to perform this request
	errCode := datastore_db.VerifyReadAccessByDataIdAndPID(types.NewUInt64(uint64(param.DataID)), connection.PID())
	if errCode != nil {
		return nil, errCode
	}

	bucket := globals.DatastoreCommon.S3Bucket
	key := fmt.Sprintf("%s/%d.bin", "letters", param.DataID)

	URL, err := globals.DatastoreCommon.S3Presigner.GetObject(bucket, key, time.Minute*15)

	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "change_error")
	}

	size, err := globals.S3ObjectSize(bucket, key)

	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "change_error")
	}

	requestHeaders, errCode := globals.DatastoreCommon.S3PostRequestHeaders()
	if errCode != nil {
		return nil, errCode
	}

	pReqPostInfo := datastore_types.NewDataStoreReqGetInfoV1()

	pReqPostInfo.URL = types.NewString(URL.String())
	pReqPostInfo.RootCACert = types.NewBuffer(globals.DatastoreCommon.RootCACert)
	pReqPostInfo.RequestHeaders = requestHeaders
	pReqPostInfo.Size = types.NewUInt32(uint32(size))

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	pReqPostInfo.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodPrepareGetObjectV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
