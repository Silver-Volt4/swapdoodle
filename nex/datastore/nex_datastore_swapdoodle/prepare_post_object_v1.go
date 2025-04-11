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

	// TODO - Need to verify what param.PersistenceInitParam.DeleteLastObject really means. It's often set to true even when it wouldn't make sense
	dataID, errCode := datastore_db.InitializeObjectByPreparePostParam(connection.PID(), param)
	if errCode != nil {
		globals.Logger.Errorf("Error on object init: %s", errCode.Error())
		return nil, errCode
	}

	bucket := globals.DatastoreCommon.S3Bucket
	key := fmt.Sprintf("%s/%d.bin", "letters", dataID)

	URL, formData, err := globals.DatastoreCommon.S3Presigner.PostObject(bucket, key, time.Minute*15)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "change_error")
	}

	requestHeaders, errCode := globals.DatastoreCommon.S3PostRequestHeaders()
	if errCode != nil {
		return nil, errCode
	}

	pReqPostInfo := datastore_types.NewDataStoreReqPostInfoV1()

	pReqPostInfo.DataID = types.NewUInt32(dataID)
	pReqPostInfo.URL = types.NewString(URL.String())
	pReqPostInfo.RequestHeaders = types.NewList[datastore_types.DataStoreKeyValue]()
	pReqPostInfo.FormFields = types.NewList[datastore_types.DataStoreKeyValue]()
	pReqPostInfo.RootCACert = types.NewBuffer(globals.DatastoreCommon.RootCACert)
	pReqPostInfo.RequestHeaders = requestHeaders

	for key, value := range formData {
		field := datastore_types.NewDataStoreKeyValue()
		field.Key = types.NewString(key)
		field.Value = types.NewString(value)

		pReqPostInfo.FormFields = append(pReqPostInfo.FormFields, field)
	}

	rmcResponseStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	pReqPostInfo.WriteTo(rmcResponseStream)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCSuccess(endpoint, rmcResponseBody)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodPreparePostObjectV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
