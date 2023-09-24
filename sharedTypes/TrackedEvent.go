package sharedTypes

type TrackedEvent struct {
	ID          string
	Camera      string
	Zones       []string
	Label       string
	HasSnapshot bool
}
