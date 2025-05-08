package webhook

import (
	"encoding/json"
	"healthmonitor/internal/db"
	"net/http"

	"github.com/rs/zerolog/log"
)

type WebhookHandler struct {
	dbh *db.DB
}

func New(dbh *db.DB) (*WebhookHandler, error) {
	return &WebhookHandler{
		dbh: dbh,
	}, nil
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var notification PlatformNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if notification.EventTopic != "device" {
		// Just discard the notification if it's not related to a device
		w.WriteHeader(http.StatusOK)
		return
	}

	data := notification.EventData

	err := h.dbh.UpsertDevice(
		data.DeviceID, data.DeviceName, data.DeviceDescription,
		data.SiteID, data.SiteName, data.SiteDescription)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Error().Err(err).Msg("Failed to upsert device")
		return
	}

	switch notification.EventCode {
	case "disconnected":
		err = h.dbh.SetDeviceState(data.DeviceID, "disconnected", true)
	case "connected":
		err = h.dbh.SetDeviceState(data.DeviceID, "disconnected", false)
	case "storage_disruption":
		err = h.dbh.SetDeviceState(data.DeviceID, "storage_disrupted", true)
	case "storage_healthy":
		err = h.dbh.SetDeviceState(data.DeviceID, "storage_disrupted", false)
	case "cpu_overutilization":
		err = h.dbh.SetDeviceState(data.DeviceID, "cpu_overutilized", true)
	case "cpu_overutilization_end":
		err = h.dbh.SetDeviceState(data.DeviceID, "cpu_overutilized", false)
	case "ram_overutilization":
		err = h.dbh.SetDeviceState(data.DeviceID, "ram_overutilized", true)
	case "ram_overutilization_end":
		err = h.dbh.SetDeviceState(data.DeviceID, "ram_overutilized", false)
	case "storage_full":
		err = h.dbh.SetDeviceState(data.DeviceID, "storage_full", true)
	case "storage_full_end":
		err = h.dbh.SetDeviceState(data.DeviceID, "storage_full", false)
	case "network_packetloss":
		err = h.dbh.SetDeviceState(data.DeviceID, "network_packet_loss", true)
	case "network_packetloss_end":
		err = h.dbh.SetDeviceState(data.DeviceID, "network_packet_loss", false)
	case "axis_imagehealth":
		err = h.dbh.SetDeviceState(data.DeviceID, "image_health_impaired", true)
	case "axis_imagehealth_end":
		err = h.dbh.SetDeviceState(data.DeviceID, "image_health_impaired", false)
	case "temperature_warning":
		err = h.dbh.SetDeviceState(data.DeviceID, "temperature_warning", true)
	case "temperature_warning_end":
		err = h.dbh.SetDeviceState(data.DeviceID, "temperature_warning", false)
	default:
		err = nil
	}

	if err != nil {
		log.Error().Err(err).Msg("Failed to update device status")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Info().Str("code", notification.EventCode).Msg("Notification processed successfully")

	w.WriteHeader(http.StatusOK)
}
