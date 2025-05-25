package datastore_db

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/silver-volt4/swapdoodle/database"
)

func GetSpecificMetaByIDs(pid types.PID, dataids types.List[types.UInt32]) (types.List[datastore_types.DataStoreSpecificMetaInfoV1], *nex.Error) {

	var ids string
	for i, number := range dataids {
		if i != 0 {
			ids += ","
		}
		ids += number.String()
	}

	fmt.Println(ids)

	rows, err := database.Postgres.Query(`SELECT
		data_id,
		owner,
		size,
		data_type,
		1
	FROM datastore.objects
	WHERE data_id IN (`+ids+`) AND (owner = $1 OR $1 IN (SELECT UNNEST(permission_recipients)))`, pid)

	if err != nil {
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, err.Error())
	}

	var lenDataIds int = len(dataids)

	res := make(types.List[datastore_types.DataStoreSpecificMetaInfoV1], lenDataIds)

	for i := 0; i < lenDataIds; i++ {
		if !rows.Next() {
			// One of the data IDs in the query above wasn't returned (it was filtered out by the owner = ... clause)
			return nil, nex.NewError(nex.ResultCodes.DataStore.OperationNotAllowed, "Access denied")
		}
		rows.Scan(
			&res[i].DataID,
			&res[i].OwnerID,
			&res[i].Size,
			&res[i].DataType,
			&res[i].Version,
		)
	}

	return res, nil
}
