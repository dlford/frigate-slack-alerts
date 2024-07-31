package main

import (
	"app/config"
	"app/events"
	"fmt"
	"sync"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup
	events.Start(
		&wg,
		cfg.MQTT_Broker,
		cfg.MQTT_Port,
		cfg.MQTT_User,
		cfg.MQTT_Password,
		cfg.MQTT_Client_ID,
		cfg.Frigate_Topic_Prefix,
		cfg.Frigate_Internal_BaseURL,
		cfg.Frigate_External_BaseURL,
		cfg.Slack_Token,
		cfg.Slack_Channel_ID,
		cfg.Filters,
	)
	wg.Wait()
}
