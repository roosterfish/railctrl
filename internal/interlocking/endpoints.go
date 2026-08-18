package interlocking

import (
	"net/http"

	"github.com/roosterfish/railctrl/api"
)

func getRouteFind(d *Daemon, r *http.Request) (*api.Route, error) {
	query := r.URL.Query()

	start := query.Get("start")
	end := query.Get("end")

	route, err := d.layout.Dijkstra(start, end)
	if err != nil {
		return nil, err
	}

	return route, nil
}
