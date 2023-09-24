package alert

import (
	"app/sharedTypes"
	"fmt"
	"io"
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

	slackClient := slack.New(slackToken)

	if !values.HasSnapshot {
		slackClient.PostMessage(slackChannel, slack.MsgOptionText(messageText, false))
	} else {
		imgURL := fmt.Sprintf("%s/api/events/%s/snapshot.jpg", internalBaseURL, values.ID)
		imgFile := fmt.Sprintf("tmp/%s.jpg", values.ID)

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

		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		fileUploadParams := slack.FileUploadParameters{
			Filetype:       "image/jpg",
			Filename:       "snapshot.jpg",
			Channels:       []string{slackChannel},
			InitialComment: messageText,
			File:           imgFile,
		}
		_, err = slackClient.UploadFile(fileUploadParams)
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
