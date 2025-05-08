package migrations

import (
	"context"
	"database/sql"
)

func Up00001(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `CREATE TABLE devices (
		id TEXT PRIMARY KEY NOT NULL,
		name TEXT NOT NULL,
		description TEXT,
		site_id TEXT,
		site_name TEXT,
		site_description TEXT,
		disconnected boolean DEFAULT false,
		storage_disrupted boolean DEFAULT false,
		cpu_overutilized boolean DEFAULT false,
		ram_overutilized boolean DEFAULT false,
		storage_full boolean DEFAULT false,
		network_packet_loss boolean DEFAULT false,
		image_health_impaired boolean DEFAULT false
	);`); err != nil {
		return err
	}

	return nil
}

func Down00001(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "DROP TABLE devices;"); err != nil {
		return err
	}

	return nil
}
