package events

import (
	"app/alert"
	"app/filter"
	"app/sharedTypes"

	"encoding/json"
	"fmt"
	_ "image/jpeg"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Start(
	wg *sync.WaitGroup,
	broker string,
	port int,
	user string,
	password string,
	clientID string,
	topicPrefix string,
	internalBaseURL string,
	externalBaseURL string,
	slackToken string,
	slackChannel string,
	filters sharedTypes.FilterConfig,
) {
	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		fmt.Println("Connected to mqtt broker")
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		fmt.Printf("Connect to mqtt broker lost: %v\n", err)
		panic(err)
	}

	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		if msg.Retained() {
			return
		}

		var message sharedTypes.EventMessage
		err := json.Unmarshal(msg.Payload(), &message)
		if err != nil {
			fmt.Println(err)
		}

		values, shouldSend := filter.FilterMessage(message, filters)
		if shouldSend {
			alert.SendAlert(
				slackChannel,
				slackToken,
				values,
				internalBaseURL,
				externalBaseURL,
			)
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	if user != "" && password != "" {
		opts.SetUsername(user)
		opts.SetPassword(password)
	}
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	client.Subscribe(fmt.Sprintf("%s/events", topicPrefix), 0, nil)

	wg.Add(1)
}
