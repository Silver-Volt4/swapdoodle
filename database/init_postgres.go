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

	_, err = Postgres.Exec(`CREATE SEQUENCE IF NOT EXISTS datastore.object_data_id_seq
		INCREMENT 1
		MINVALUE 1
		MAXVALUE 281474976710656
		START 1
		CACHE 1`,
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.objects (
		data_id BIGINT NOT NULL DEFAULT nextval('datastore.object_data_id_seq') PRIMARY KEY,
		upload_completed BOOLEAN NOT NULL DEFAULT FALSE,
		deleted BOOLEAN NOT NULL DEFAULT FALSE,
		owner INT,
		size INT,
		name TEXT,
		data_type INT,
		meta_binary BYTEA,
		permission INT,
		permission_recipients INT[],
		delete_permission INT,
		delete_permission_recipients INT[],
		flag INT,
		period INT,
		refer_data_id BIGINT,
		tags TEXT[],
		access_password BIGINT NOT NULL DEFAULT 0,
		update_password BIGINT NOT NULL DEFAULT 0,
		creation_date TIMESTAMP,
		update_date TIMESTAMP
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	// In our testing, having the notification ID set to low numbers caused some weird issues.
	// Bumping it up seems to have fixed them, so let's start with a high number.
	// Since notifications are sequential, this can also remedy some issues with the order not matching up (from the time when SD still worked)
	_, err = Postgres.Exec(`CREATE SEQUENCE IF NOT EXISTS datastore.notification_id_seq
		INCREMENT 1
		START 100000000
		CACHE 1`,
	)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	_, err = Postgres.Exec(`CREATE TABLE IF NOT EXISTS datastore.notifications (
		notification_id BIGINT NOT NULL DEFAULT nextval('datastore.notification_id_seq') PRIMARY KEY,
		recipient_pid INT,
		data_id BIGINT REFERENCES datastore.objects (data_id)
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		os.Exit(0)
	}

	globals.Logger.Success("Postgres tables created")
}
