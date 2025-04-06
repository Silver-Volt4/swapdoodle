package datastore_smm_db

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/silver-volt4/swapdoodle/database"
	datastore_db "github.com/silver-volt4/swapdoodle/database/datastore"
	"github.com/silver-volt4/swapdoodle/globals"
)

func InsertOrUpdateCustomRanking(dataID uint64, applicationID, score uint32) uint32 {
	errCode := datastore_db.IsObjectAvailable(dataID)
	if errCode != 0 {
		globals.Logger.Errorf("Error code %d", errCode)
		return errCode
	}

	_, err := database.Postgres.Exec(`INSERT INTO datastore.object_custom_rankings (
		data_id,
		application_id,
		value
	) VALUES (
		$1,
		$2,
		$3
	) ON CONFLICT (data_id, application_id) DO UPDATE SET value=datastore.object_custom_rankings.value+EXCLUDED.value`,
		dataID,
		applicationID,
		score,
	)

	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.Errors.DataStore.Unknown
	}

	return 0
}
