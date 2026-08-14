package api

type TrackType string

type Turnout struct {
	ThrowID           int `yaml:"throw_id"`
	ReversePolarityID int `yaml:"reverse_polarity_id"`
}

type Track struct {
	ID          string    `yaml:"id"`
	Type        TrackType `yaml:"type"`
	Length      string    `yaml:"length"`
	OccupancyID int       `yaml:"occupancy_id"`

	// Use pointer to make it an optional track property.
	*Turnout `yaml:",omitempty,inline"`
}

type Connection string
type Connections map[string][]Connection

type SignalType string
type SignalAspect struct {
	Name string `yaml:"name"`
}

type Signal struct {
	ID       string         `yaml:"id"`
	From     string         `yaml:"from"`
	To       string         `yaml:"to"`
	Type     SignalType     `yaml:"type"`
	Distance string         `yaml:"distance"`
	Aspects  []SignalAspect `yaml:"aspects"`
}

type Layout struct {
	Version     string      `yaml:"version"`
	Tracks      []Track     `yaml:"tracks"`
	Connections Connections `yaml:"connections"`
	Signals     []Signal    `yaml:"signals"`
}
