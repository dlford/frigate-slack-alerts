package filter

import (
	"app/sharedTypes"
	"fmt"
)

var trackedEvents = make(map[string]sharedTypes.TrackedEvent)

// TODO: Send initial message to slack and update as things change (https://pkg.go.dev/github.com/slack-go/slack@v0.12.3#Client.UpdateMessage)

func FilterMessage(message sharedTypes.EventMessage, config sharedTypes.FilterConfig) (sharedTypes.TrackedEvent, bool) {
	tracked := trackedEvents[message.After.ID]

	tracked.ID = message.After.ID
	tracked.Camera = message.After.Camera
	tracked.Label = message.After.Label
	tracked.HasSnapshot = message.After.HasSnapshot
	for cz := range message.After.CurrentZones {
		found := false
		for tz := range tracked.Zones {
			if tracked.Zones[tz] == message.After.CurrentZones[cz] {
				found = true
				break
			}
		}
		if !found {
			tracked.Zones = append(tracked.Zones, message.After.CurrentZones[cz])
		}
	}

	trackedEvents[message.After.ID] = tracked

	fmt.Printf("tracking: %+v\n", tracked)

	if message.Type != "end" {
		return tracked, false
	}

	if len(config.Cameras) <= 0 {
		fmt.Printf("alerted:  %+v\n", tracked)
		return tracked, true
	}

	// Camera
	for c := range config.Cameras {
		if config.Cameras[c].Name == tracked.Camera {
			// Zone
			for z := range config.Cameras[c].Zones {
				for cz := range tracked.Zones {
					if config.Cameras[c].Zones[z].Name == tracked.Zones[cz] {
						// Object
						for o := range config.Cameras[c].Zones[z].Objects {
							if config.Cameras[c].Zones[z].Objects[o] == tracked.Label {
								fmt.Printf("alerted:  %+v\n", tracked)
								delete(trackedEvents, message.After.ID)
								return tracked, true
							}
						}
					}
				}
			}
			break
		}
	}

	fmt.Printf("ignored:  %+v\n", tracked)

	delete(trackedEvents, message.After.ID)
	return tracked, false
}
