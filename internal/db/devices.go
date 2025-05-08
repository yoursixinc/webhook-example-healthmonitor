package db

import (
	"fmt"
)

type Device struct {
	Id                  string
	Name                string
	Description         string
	SiteId              string
	SiteName            string
	SiteDescription     string
	Disconnected        bool
	StorageDisrupted    bool
	CpuOverutilized     bool
	RamOverutilized     bool
	StorageFull         bool
	NetworkPacketLoss   bool
	ImageHealthImpaired bool
}

type DeviceCounts struct {
	Tracked             int
	Disconnected        int
	StorageDisrupted    int
	CpuOverutilized     int
	RamOverutilized     int
	StorageFull         int
	NetworkPacketLoss   int
	ImageHealthImpaired int
}

func (db *DB) UpsertDevice(id, name, description, siteId, siteName, siteDescription string) error {
	query := `INSERT INTO devices(id, name, description, site_id, site_name, site_description) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET name = $2, description = $3, site_id = $4, site_name = $5, site_description = $6;`

	_, err := db.sqldb.Exec(query,
		id,
		name,
		description,
		siteId,
		siteName,
		siteDescription,
	)

	return err
}

func (db *DB) SetDeviceState(id, stateName string, stateValue bool) error {
	query := fmt.Sprintf(`UPDATE devices SET %s = $1 WHERE id = $2`, stateName)

	_, err := db.sqldb.Exec(query,
		stateValue,
		id,
	)
	return err
}

func (db *DB) GetDevices() ([]Device, error) {
	query := `SELECT id, name, description, site_id, site_name, site_description,
		disconnected, storage_disrupted, cpu_overutilized, ram_overutilized,
		storage_full, network_packet_loss, image_health_impaired FROM devices;`

	rows, err := db.sqldb.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	devices := []Device{}
	for rows.Next() {
		var device Device
		if err := rows.Scan(
			&device.Id,
			&device.Name,
			&device.Description,
			&device.SiteId,
			&device.SiteName,
			&device.SiteDescription,
			&device.Disconnected,
			&device.StorageDisrupted,
			&device.CpuOverutilized,
			&device.RamOverutilized,
			&device.StorageFull,
			&device.NetworkPacketLoss,
			&device.ImageHealthImpaired); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
}

func (db *DB) GetDeviceCounts() (*DeviceCounts, error) {
	query := `SELECT
		COALESCE(COUNT(*), 0) AS tracked,
		COALESCE(SUM(CASE WHEN disconnected THEN 1 ELSE 0 END), 0) AS disconnected,
		COALESCE(SUM(CASE WHEN storage_disrupted THEN 1 ELSE 0 END), 0) AS storage_disrupted,
		COALESCE(SUM(CASE WHEN cpu_overutilized THEN 1 ELSE 0 END), 0) AS cpu_overutilized,
		COALESCE(SUM(CASE WHEN ram_overutilized THEN 1 ELSE 0 END), 0) AS ram_overutilized,
		COALESCE(SUM(CASE WHEN storage_full THEN 1 ELSE 0 END), 0) AS storage_full,
		COALESCE(SUM(CASE WHEN network_packet_loss THEN 1 ELSE 0 END), 0) AS network_packet_loss,
		COALESCE(SUM(CASE WHEN image_health_impaired THEN 1 ELSE 0 END), 0) AS image_health_impaired
	FROM devices;`

	row := db.sqldb.QueryRow(query)

	counts := &DeviceCounts{}
	if err := row.Scan(
		&counts.Tracked,
		&counts.Disconnected,
		&counts.StorageDisrupted,
		&counts.CpuOverutilized,
		&counts.RamOverutilized,
		&counts.StorageFull,
		&counts.NetworkPacketLoss,
		&counts.ImageHealthImpaired); err != nil {
		return nil, err
	}

	return counts, nil
}
