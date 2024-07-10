package alert

import (
	"app/sharedTypes"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TODO: Add buttons to link to external url paths

func SendAlert(
	slackChannel string,
	slackToken string,
	values sharedTypes.TrackedEvent,
	internalBaseURL string,
	extenalBaseURL string,
) {
	caser := cases.Title(language.English)

	messageText := fmt.Sprintf(
		"%s seen near %s",
		caser.String(strings.Replace(values.Label, "_", " ", -1)),
		caser.String(strings.Replace(values.Camera, "_", " ", -1)),
	)

	if len(values.Zones) > 0 {
		var zones []string
		for z := range values.Zones {
			zones = append(zones, caser.String(strings.Replace(values.Zones[z], "_", " ", -1)))
		}
		messageText = fmt.Sprintf("%s in zones: %s.", messageText, strings.Join(zones, ", "))
	}

	messageText = fmt.Sprintf("%s\nView events: %s/events", messageText, extenalBaseURL)

	slackClient := slack.New(slackToken)

	if !values.HasSnapshot {
		slackClient.PostMessage(slackChannel, slack.MsgOptionText(messageText, false))
	} else {
		imgURL := fmt.Sprintf("%s/api/events/%s/snapshot.jpg", internalBaseURL, values.ID)
		imgFile := fmt.Sprintf("/tmp/%s.jpg", values.ID)

		response, err := http.Get(imgURL)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer response.Body.Close()

		file, err := os.Create(imgFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		size, err := io.Copy(file, response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		if size > math.MaxInt32 {
			fmt.Println("Snapshot file is too large to upload")
			return
		}

		// Ideally, the cursor should be used to support teams with >200 channels
		conversationsParams := slack.GetConversationsParameters{
			Limit: 200,
		}

		channels, _, err := slackClient.GetConversations(&conversationsParams)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err != nil {
			fmt.Println(err)
		}

		channelId := ""
		for _, channel := range channels {
			if slackChannel == channel.Name {
				channelId = channel.GroupConversation.Conversation.ID
				break
			}
		}

		if channelId == "" {
			fmt.Println("Couldn't find slack channel")
		}

		fileUploadV2Params := slack.UploadFileV2Parameters{
			File:           imgFile,
			FileSize:       int(size),
			Filename:       "snapshot.jpg",
			InitialComment: messageText,
			Channel: 	    channelId,
		}

		_, err = slackClient.UploadFileV2(fileUploadV2Params)
		if err != nil {
			fmt.Println(err)
			return
		}

		if imgFile != "" {
			if err := os.Remove(imgFile); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
