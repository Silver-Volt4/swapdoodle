package datastore_db

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

func DeleteObjectByDataID(dataID uint64) uint32 {
	errCode := IsObjectAvailable(dataID)
	if errCode != 0 {
		return errCode
	}

	_, err := database.Postgres.Exec(`UPDATE datastore.objects SET deleted=TRUE WHERE data_id=$1`, dataID)
	if err != nil {
		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.Errors.DataStore.Unknown
	}

	return 0
}
