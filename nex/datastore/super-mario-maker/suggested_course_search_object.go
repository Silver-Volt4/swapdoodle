package nex_datastore_super_mario_maker

import (
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	datastore_smm_db "github.com/silver-volt4/swapdoodle/database/datastore/super-mario-maker"
	"github.com/silver-volt4/swapdoodle/globals"
)

func SuggestedCourseSearchObject(err error, packet nex.PacketInterface, callID uint32, param datastore_types.DataStoreSearchParam, extraData types.List[types.String]) (*nex.RMCMessage, *nex.Error) {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// * This method is called when a course is completed
	// * to show the scrolling courses at the bottom of the
	// * screen. extraData[0] is the DataID of the current
	// * course. extraData[1] and extraData[2] are usually
	// * 2 and 6 respectively. extraData[3] and extraData[4]
	// * are always 0? extraData seems to not have any
	// * effect on the NUMBER of courses returned but likely
	// * does act as a filter of some kind? Maybe it has to
	// * do with difficulty? Or ratings?

	_, err = strconv.ParseUint(string(extraData[0]), 0, 64)
	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.InvalidArgument, "Invalid argument")
	}

	// TODO - Use extraData for filtering
	pRankingResults, nexError := datastore_smm_db.GetRandomCoursesWithLimit(int(param.ResultRange.Length))
	if nexError != nil {
		return nil, nexError
	}

	rmcResponseStream := nex.NewByteStreamOut(globals.HppServer.LibraryVersions(), globals.HppServer.ByteStreamSettings())

	pRankingResults.WriteTo(rmcResponseStream)

	rmcResponse := nex.NewRMCSuccess(globals.HppServer, rmcResponseStream.Bytes())
	rmcResponse.ProtocolID = datastore_super_mario_maker.ProtocolID
	rmcResponse.MethodID = datastore_super_mario_maker.MethodSuggestedCourseSearchObject
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
