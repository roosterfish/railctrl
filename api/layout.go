package api

import (
	"fmt"
	"slices"
)

const LayoutVersion = "1"

type TrackType string

type Turnout struct {
	ThrowID           int `yaml:"throw_id"`
	ReversePolarityID int `yaml:"reverse_polarity_id"`
}

type Track struct {
	ID          string    `yaml:"id"`
	Type        TrackType `yaml:"type"`
	Length      int       `yaml:"length"`
	OccupancyID int       `yaml:"occupancy_id"`

	// Use pointer to make it an optional track property.
	*Turnout `yaml:",omitempty,inline"`
}

type Connection struct {
	To                string `yaml:"to"`
	Branch            bool   `yaml:"branch"`
	OppositeDirection bool   `yaml:"opposite_direction"`
}

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

// Dijkstra finds the least expensive route between start and end tracks.
func (l *Layout) Dijkstra(start, end string) ([]string, int, error) {
	// Max int.
	const infinity = int(^uint(0) >> 1)

	// Look up tracks by their ID.
	tracks := make(map[string]Track, len(l.Tracks))
	for _, track := range l.Tracks {
		tracks[track.ID] = track
	}

	// Look up connections by their ID.
	connections := make(map[string][]Connection, len(l.Connections))
	for from, fromConnections := range l.Connections {
		connections[from] = fromConnections
	}

	_, ok := tracks[start]
	if !ok {
		return nil, 0, fmt.Errorf("invalid track %q", start)
	}

	_, ok = tracks[end]
	if !ok {
		return nil, 0, fmt.Errorf("invalid track %q", end)
	}

	dist := make(map[string]int)
	prev := make(map[string]string)
	visited := make(map[string]bool)

	for _, track := range l.Tracks {
		dist[track.ID] = infinity
	}

	dist[start] = 0

	for {
		// Find the unvisited track with the smallest distance.
		current := ""
		for id, d := range dist {
			if !visited[id] && (current == "" || d < dist[current]) {
				current = id
			}
		}

		// No reachable tracks left.
		if current == "" || dist[current] == infinity {
			break
		}

		visited[current] = true

		if current == end {
			break
		}

		// Visit connected tracks.
		for _, next := range l.Connections[current] {
			nextTrack := tracks[next.To]

			// In case of a turnout (>2 connections), we cannot traverse it coming from one branch and exiting the other.
			if len(l.Connections[current]) > 2 {
				ok := slices.ContainsFunc(l.Connections[current], func(c Connection) bool {
					_, ok := visited[c.To]
					return ok && c.Branch && next.Branch
				})
				if ok {
					continue
				}
			}

			weight, err := nextTrack.Weight()
			if err != nil {
				return nil, 0, fmt.Errorf("failed getting weight of track %q: %w", nextTrack.ID, err)
			}

			// Multiply the cost when the transition/direction between tracks indicates use of the opposite track.
			if next.OppositeDirection {
				weight *= 2
			}

			newDist := dist[current] + weight

			if newDist < dist[next.To] {
				dist[next.To] = newDist
				prev[next.To] = current
			}
		}
	}

	if dist[end] == infinity {
		return nil, 0, fmt.Errorf("unable to find a route from %q to %q", start, end)
	}

	// Reconstruct path.
	path := []string{}
	for at := end; at != ""; at = prev[at] {
		path = append(path, at)
	}

	// Reverse path.
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, dist[end], nil
}

// Weight returns a track's weight when passing it.
func (t *Track) Weight() (int, error) {
	return t.Length, nil
}
