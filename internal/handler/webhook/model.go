package webhook

type PlatformNotification struct {
	EventTopic       string `json:"eventTopic"`
	EventCode        string `json:"eventCode"`
	EventDescription string `json:"eventDescription"`
	EventOccurred    string `json:"eventOccurred"`
	EventData        struct {
		DeviceID          string `json:"deviceID"`
		DeviceName        string `json:"deviceName"`
		DeviceDescription string `json:"deviceDescription"`
		SiteID            string `json:"siteID"`
		SiteName          string `json:"siteName"`
		SiteDescription   string `json:"siteDescription"`
		GroupID           string `json:"groupID"`
		GroupName         string `json:"groupName"`
		GroupDescription  string `json:"groupDescription"`
		VideoSourceId     string `json:"videoSourceId"`
		DiskId            string `json:"diskId"`
		OriginCode        string `json:"originCode"`
		OriginData        struct {
			UserId         string `json:"userId"`
			UserName       string `json:"userName"`
			UserEmail      string `json:"userEmail"`
			PeripheralId   string `json:"peripheralId"`
			PeripheralName string `json:"peripheralName"`
		} `json:"originData"`
	}
}
