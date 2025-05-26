package nex_datastore_swapdoodle

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/silver-volt4/swapdoodle/globals"
)

func CompletePostObjectV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreCompletePostParamV1) (*nex.RMCMessage, *nex.Error) {
	if globals.MinIOClient == nil {
		globals.Logger.Warning("MinIOClient not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if globals.DatastoreCommon.GetObjectInfoByDataID == nil {
		globals.Logger.Warning("GetObjectInfoByDataID not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if globals.DatastoreCommon.GetObjectOwnerByDataID == nil {
		globals.Logger.Warning("GetObjectOwnerByDataID not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if globals.DatastoreCommon.GetObjectSizeByDataID == nil {
		globals.Logger.Warning("GetObjectSizeByDataID not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if globals.DatastoreCommon.UpdateObjectUploadCompletedByDataID == nil {
		globals.Logger.Warning("UpdateObjectUploadCompletedByDataID not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if globals.DatastoreCommon.DeleteObjectByDataID == nil {
		globals.Logger.Warning("DeleteObjectByDataID not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
	}

	dataID := types.NewUInt64(uint64(param.DataID))

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	// If GetObjectInfoByDataID does not return an error, it means the object has already been marked as uploaded
	_, errCode := globals.DatastoreCommon.GetObjectInfoByDataID(dataID)
	if errCode == nil {
		return nil, nex.NewError(nex.ResultCodes.DataStore.PermissionDenied, "change_error")
	}

	// Only allow the object's owner to make this request
	ownerPID, errCode := globals.DatastoreCommon.GetObjectOwnerByDataID(dataID)
	if errCode != nil {
		return nil, errCode
	}

	if ownerPID != uint32(connection.PID()) {
		return nil, nex.NewError(nex.ResultCodes.DataStore.PermissionDenied, "change_error")
	}

	bucket := globals.DatastoreCommon.S3Bucket
	key := globals.S3GetLetterKey(param.DataID)

	if param.IsSuccess {
		objectSizeS3, err := globals.DatastoreCommon.S3ObjectSize(bucket, key)
		if err != nil {
			globals.Logger.Error(err.Error())
			return nil, nex.NewError(nex.ResultCodes.DataStore.NotFound, "change_error")
		}

		objectSizeDB, nErr := globals.DatastoreCommon.GetObjectSizeByDataID(dataID)
		if nErr != nil {
			return nil, nErr
		}

		if objectSizeS3 != uint64(objectSizeDB) {
			globals.Logger.Errorf("Object with DataID %d did not upload correctly! Mismatched sizes", dataID)
			// TODO - Is this a good error?
			return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
		}

		nErr = globals.DatastoreCommon.UpdateObjectUploadCompletedByDataID(dataID, true)
		if nErr != nil {
			return nil, nErr
		}
	} else {
		nErr := globals.DatastoreCommon.DeleteObjectByDataID(dataID)
		if nErr != nil {
			return nil, nErr
		}
	}

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodCompletePostObjectV1
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
