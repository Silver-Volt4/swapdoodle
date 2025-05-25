package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

func VerifyReadAccessByDataIdAndPID(dataID types.UInt64, pid types.PID) *nex.Error {
	nexError := IsObjectAvailable(dataID)
	if nexError != nil {
		return nexError
	}

	var t bool
	err := database.Postgres.QueryRow(`SELECT 1 FROM datastore.objects WHERE data_id=$1 AND (owner = $2 OR $2 IN (SELECT UNNEST(permission_recipients)))`, dataID, pid).Scan(&t)

	if err != nil {
		if err == sql.ErrNoRows {
			return nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "Access denied")
		}

		globals.Logger.Error(err.Error())

		// TODO - Send more specific errors?
		return nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	return nil
}
