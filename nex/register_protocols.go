package nex

import (
	"os"

	"github.com/PretendoNetwork/nex-go/v2/types"
	datastorecommon "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	securecommon "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	datastore_db "github.com/silver-volt4/swapdoodle/database/datastore"
	"github.com/silver-volt4/swapdoodle/globals"
	nex_datastore_swapdoodle "github.com/silver-volt4/swapdoodle/nex/datastore/nex_datastore_swapdoodle"
)

func registerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.HppServer.RegisterServiceProtocol(secureProtocol)
	commonSecureProtocol := securecommon.NewCommonProtocol(secureProtocol)
	commonSecureProtocol.CreateReportDBRecord = func(pid types.PID, reportID types.UInt32, reportData types.QBuffer) error {
		return nil
	}

	sdDatastore := datastore.NewProtocol()

	sdDatastore.GetNotificationURL = nex_datastore_swapdoodle.GetNotificationURL
	sdDatastore.PreparePostObjectV1 = nex_datastore_swapdoodle.PreparePostObjectV1
	sdDatastore.CompletePostObjectV1 = nex_datastore_swapdoodle.CompletePostObjectV1
	sdDatastore.PrepareGetObjectV1 = nex_datastore_swapdoodle.PrepareGetObjectV1
	sdDatastore.GetSpecificMetaV1 = nex_datastore_swapdoodle.GetSpecificMetaV1
	sdDatastore.GetNewArrivedNotificationsV1 = nex_datastore_swapdoodle.GetNewArrivedNotificationsV1

	globals.HppServer.RegisterServiceProtocol(sdDatastore)

	commonDataStoreProtocol := datastorecommon.NewCommonProtocol(sdDatastore)

	commonDataStoreProtocol.SetMinIOClient(globals.MinIOClient)
	commonDataStoreProtocol.S3Bucket = os.Getenv("PN_SD_CONFIG_S3_BUCKET")
	commonDataStoreProtocol.S3Presigner = globals.Presigner
	commonDataStoreProtocol.GetObjectInfoByDataID = datastore_db.GetObjectInfoByDataID
	commonDataStoreProtocol.GetObjectOwnerByDataID = datastore_db.GetObjectOwnerByDataID
	commonDataStoreProtocol.GetObjectSizeByDataID = datastore_db.GetObjectSizeByDataID
	commonDataStoreProtocol.UpdateObjectUploadCompletedByDataID = datastore_db.UpdateObjectUploadCompletedByDataID
	commonDataStoreProtocol.DeleteObjectByDataID = datastore_db.DeleteObjectByDataID

	globals.DatastoreCommon = commonDataStoreProtocol
}
