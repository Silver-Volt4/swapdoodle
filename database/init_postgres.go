package database

import (
	"os"

	"github.com/silver-volt4/swapdoodle/globals"
)

func initPostgres() {
	_, err := Postgres.Exec(`CREATE SCHEMA IF NOT EXISTS datastore`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("datastore Postgres schema created")

	// * Super Mario Maker has non-standard DataID requirements.
	// * DataID 900000 is reserved for the event course metadata
	// * file, and official event courses begin at DataID 930010
	// * and end at DataID 930050. To prevent a collision
	// * eventually, we need to start course IDs AFTER 930050
	// *
	// * DataIDs are stored and processed as uint64, however
	// * Super Mario Maker can not use the full uint64 range.
	// * This is because course share codes are generated from the
	// * courses DataID. A course share code is an 8 byte hex
	// * string, where the upper 2 bytes are the checksum of the
	// * lower 6 bytes. The lower 6 bytes are the courses DataID
	// *
	// * Super Mario Maker is only capable of displaying codes up
	// * to 0xFFFFFFFFFFFF, essentially truncating DataIDs down to
	// * 48 bit integers instead of 64 bit. I doubt we will ever
	// * hit even the 32 bit limit, let alone 48, but this is here
	// * just in case
	_, err = Postgres.Exec(`CREATE SEQUENCE IF NOT EXISTS datastore.object_data_id_seq
		INCREMENT 1
		MINVALUE 1
		MAXVALUE 281474976710656
		START 940000
		CACHE 1`, // * Honestly I don't know what CACHE does but I saw it recommended so here it is
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.objects (
		data_id bigint NOT NULL DEFAULT nextval('datastore.object_data_id_seq') PRIMARY KEY,
		upload_completed boolean NOT NULL DEFAULT FALSE,
		deleted boolean NOT NULL DEFAULT FALSE,
		owner int,
		size int,
		name text,
		data_type int,
		meta_binary bytea,
		permission int,
		permission_recipients int[],
		delete_permission int,
		delete_permission_recipients int[],
		flag int,
		period int,
		refer_data_id bigint,
		tags text[],
		access_password bigint NOT NULL DEFAULT 0,
		update_password bigint NOT NULL DEFAULT 0,
		creation_date timestamp,
		update_date timestamp
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("Postgres tables created")
}
