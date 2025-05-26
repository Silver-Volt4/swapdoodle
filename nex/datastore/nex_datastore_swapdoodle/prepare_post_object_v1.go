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

func PreparePostObjectV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStorePreparePostParamV1) (*nex.RMCMessage, *nex.Error) {
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

	_dataID, nErr := datastore_db.InitializeObjectByPreparePostParam(connection.PID(), param)
	dataID := types.NewUInt32(_dataID)

	if nErr != nil {
		globals.Logger.Errorf("Error on object init: %s", nErr.Error())
		return nil, nErr
	}

	bucket := globals.DatastoreCommon.S3Bucket
	key := globals.S3GetLetterKey(dataID)

	URL, formData, err := globals.DatastoreCommon.S3Presigner.PostObject(bucket, key, time.Minute*15)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "change_error")
	}

	requestHeaders, nErr := globals.DatastoreCommon.S3PostRequestHeaders()
	if nErr != nil {
		return nil, nErr
	}

	resStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	postObjectInfo := datastore_types.NewDataStoreReqPostInfoV1()
	postObjectInfo.DataID = dataID
	postObjectInfo.URL = types.NewString(URL.String())
	postObjectInfo.RequestHeaders = types.NewList[datastore_types.DataStoreKeyValue]()
	postObjectInfo.FormFields = types.NewList[datastore_types.DataStoreKeyValue]()
	postObjectInfo.RootCACert = types.NewBuffer(globals.DatastoreCommon.RootCACert)
	postObjectInfo.RequestHeaders = requestHeaders

	for key, value := range formData {
		field := datastore_types.NewDataStoreKeyValue()
		field.Key = types.NewString(key)
		field.Value = types.NewString(value)
		postObjectInfo.FormFields = append(postObjectInfo.FormFields, field)
	}

	postObjectInfo.WriteTo(resStream)

	res := nex.NewRMCSuccess(endpoint, resStream.Bytes())
	res.ProtocolID = datastore.ProtocolID
	res.MethodID = datastore.MethodPreparePostObjectV1
	res.CallID = callID

	return res, nil
}
