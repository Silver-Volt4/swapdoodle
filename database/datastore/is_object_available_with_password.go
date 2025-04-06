package datastore_db

import (
	"database/sql"

	"github.com/PretendoNetwork/nex-go"
	"github.com/silver-volt4/swapdoodle/database"
	"github.com/silver-volt4/swapdoodle/globals"
)

func IsObjectAvailableWithPassword(dataID, password uint64) uint32 {
	var underReview bool
	var accessPassword uint64

	err := database.Postgres.QueryRow(`SELECT
		under_review,
		access_password
	FROM datastore.objects WHERE data_id=$1 AND upload_completed=TRUE AND deleted=FALSE`, dataID).Scan(
		&underReview,
		&accessPassword,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nex.Errors.DataStore.NotFound
		}

		globals.Logger.Error(err.Error())
		// TODO - Send more specific errors?
		return nex.Errors.DataStore.Unknown
	}

	if accessPassword != 0 && accessPassword != password {
		return nex.Errors.DataStore.InvalidPassword
	}

	if underReview {
		return nex.Errors.DataStore.UnderReviewing
	}

	return 0
}
