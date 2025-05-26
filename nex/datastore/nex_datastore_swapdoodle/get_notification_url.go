package nex_datastore_swapdoodle

import (
	"fmt"
	"strings"
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
	key := globals.S3GetNotificationKey(packet.Sender().PID())

	url, err := globals.DatastoreCommon.S3Presigner.GetObject(bucket, key, time.Hour*24*7)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "change_error")
	}

	resStream := nex.NewByteStreamOut(globals.HppServer.LibraryVersions(), globals.HppServer.ByteStreamSettings())

	urlInfo := datastore_types.NewDataStoreReqGetNotificationURLInfo()

	urlInfo.URL = types.NewString(fmt.Sprintf("%s://%s/", url.Scheme, url.Host))
	urlInfo.Key = types.NewString(strings.TrimPrefix(url.Path, "/"))
	urlInfo.Query = types.NewString(fmt.Sprintf("?%s", url.Query().Encode()))
	urlInfo.RootCACert = types.NewBuffer(nil)

	urlInfo.WriteTo(resStream)

	res := nex.NewRMCSuccess(globals.HppServer, resStream.Bytes())
	res.ProtocolID = datastore.ProtocolID
	res.MethodID = datastore.MethodGetNotificationURL
	res.CallID = callID

	return res, nil
}
