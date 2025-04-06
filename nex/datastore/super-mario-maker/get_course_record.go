package nex_datastore_super_mario_maker

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_super_mario_maker_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker/types"
	datastore_smm_db "github.com/silver-volt4/swapdoodle/database/datastore/super-mario-maker"
	"github.com/silver-volt4/swapdoodle/globals"
)

func GetCourseRecord(err error, packet nex.PacketInterface, callID uint32, param datastore_super_mario_maker_types.DataStoreGetCourseRecordParam) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	result, nexError := datastore_smm_db.GetCourseRecordByDataIDAndSlot(param.DataID, param.Slot)
	if nexError != nil {
		return nil, nexError
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.HppServer.LibraryVersions(), globals.HppServer.ByteStreamSettings())

	result.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.HppServer, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodGetCourseRecord
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
