package datastore_db

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

func UpdateObjectUploadCompletedByDataID(dataID types.UInt64, uploadCompleted bool) *nex.Error {
	_, err := database.Postgres.Exec(`UPDATE datastore.objects SET upload_completed=$1 WHERE data_id=$2`, uploadCompleted, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// Create notifications
	_, err = database.Postgres.Exec(`INSERT INTO datastore.notifications (data_id, recipient_pid)
		SELECT data_id, UNNEST(permission_recipients) as recipient_pid
		FROM datastore.objects
		WHERE data_id = $1`, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
