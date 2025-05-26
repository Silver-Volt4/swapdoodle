package nex_datastore_swapdoodle

import (
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

	// Only allow the owner or recipient to perform this request
	nErr := datastore_db.VerifyReadAccessByDataIdAndPID(types.NewUInt64(uint64(param.DataID)), connection.PID())
	if nErr != nil {
		return nil, nErr
	}

	bucket := globals.DatastoreCommon.S3Bucket
	key := globals.S3GetLetterKey(param.DataID)

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

	requestHeaders, nErr := globals.DatastoreCommon.S3PostRequestHeaders()
	if nErr != nil {
		return nil, nErr
	}

	resStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	getObjectInfo := datastore_types.NewDataStoreReqGetInfoV1()
	getObjectInfo.URL = types.NewString(URL.String())
	getObjectInfo.RootCACert = types.NewBuffer(globals.DatastoreCommon.RootCACert)
	getObjectInfo.RequestHeaders = requestHeaders
	getObjectInfo.Size = types.NewUInt32(uint32(size))

	getObjectInfo.WriteTo(resStream)

	res := nex.NewRMCSuccess(endpoint, resStream.Bytes())
	res.ProtocolID = datastore.ProtocolID
	res.MethodID = datastore.MethodPrepareGetObjectV1
	res.CallID = callID

	return res, nil
}
