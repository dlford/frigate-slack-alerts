package sharedTypes

type FilterZone struct {
	Name    string   `yaml:"name"`
	Objects []string `yaml:"objects"`
}

type FilterCamera struct {
	Name  string       `yaml:"name"`
	Zones []FilterZone `yaml:"zones"`
}

type FilterConfig struct {
	Cameras []FilterCamera `yaml:"cameras"`
}

type Config struct {
	MQTT_User                      string
	MQTT_Password                  string
	MQTT_Broker                    string `default:"localhost"`
	MQTT_Port                      int    `default:"1883"`
	MQTT_Client_ID                 string `default:"go-frigate-slack-alerts"`
	Frigate_Internal_BaseURL       string `required:"true"`
	Frigate_External_BaseURL       string
	Frigate_Topic_Prefix           string `default:"frigate"`
	Slack_Token                    string `required:"true"`
	Slack_Channel_ID               string `required:"true"`
	Filter_Config_File             string
	Ignore_Events_Without_Snapshot bool `default:"false"`
	Filters                        FilterConfig
}
