package datastore_db

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/silver-volt4/swapdoodle/database"
)

func GetArrivedNotificationsByPID(pid types.PID, lastNotification types.UInt64, limit types.UInt16) (types.List[datastore_types.DataStoreNotificationV1], types.Bool, *nex.Error) {
	rows, err := database.Postgres.Query(`SELECT
		notification_id,
		data_id
	FROM datastore.notifications 
	WHERE recipient_pid=$1 AND notification_id>$2
	ORDER BY notification_id ASC 
	LIMIT ($3+1)`, pid, lastNotification, limit)

	if err != nil {
		return nil, false, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	res := make(types.List[datastore_types.DataStoreNotificationV1], limit)

	var i types.UInt16 = 0
	more := false

	for rows.Next() {
		if i == limit {
			more = true
			break
		}
		rows.Scan(
			&res[i].NotificationID,
			&res[i].DataID,
		)
		i++
	}

	return res[0:i], types.NewBool(more), nil
}
