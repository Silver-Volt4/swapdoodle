package datastore_db

import (
	"fmt"

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

	rows, err := database.Postgres.Query(`SELECT MAX(notification_id), recipient_pid
	FROM datastore.notifications
	WHERE recipient_pid IN (SELECT UNNEST(permission_recipients) FROM datastore.objects WHERE data_id = $1)
	GROUP BY recipient_pid`, dataID)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	// Update "noti-files"
	for rows.Next() {
		var notificationId uint64
		var pid types.PID
		rows.Scan(&notificationId, &pid)

		bucket := globals.DatastoreCommon.S3Bucket
		key := fmt.Sprintf("%s/%d", "notifications", pid)

		// TODO: Looks like we're not the only ones wondering about the meaning of that last number:
		// https://github.com/PretendoNetwork/pokemon-rumble-world-secure/blob/main/nex/datastore/complete_post_object_v1.go#L46
		// Let's leave it constant for now...
		globals.S3SetFileContent(bucket, key, fmt.Sprintf("%d,%d,%d", notificationId, pid, 1479103557))
	}

	return nil
}
