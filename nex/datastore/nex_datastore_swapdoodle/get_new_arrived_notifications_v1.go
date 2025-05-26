package nex_datastore_swapdoodle

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_db "github.com/silver-volt4/swapdoodle/database/datastore"
	"github.com/silver-volt4/swapdoodle/globals"
)

func GetNewArrivedNotificationsV1(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreGetNewArrivedNotificationsParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	notifications, more, nErr := datastore_db.GetArrivedNotificationsByPID(connection.PID(), param.LastNotificationID, param.Limit)
	if nErr != nil {
		return nil, nErr
	}

	resStream := nex.NewByteStreamOut(endpoint.LibraryVersions(), endpoint.ByteStreamSettings())

	notifications.WriteTo(resStream)
	more.WriteTo(resStream)

	res := nex.NewRMCSuccess(endpoint, resStream.Bytes())
	res.ProtocolID = datastore.ProtocolID
	res.MethodID = datastore.MethodGetNewArrivedNotificationsV1
	res.CallID = callID

	return res, nil
}
