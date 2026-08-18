package interlocking

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/roosterfish/railctrl/api"
	"github.com/roosterfish/railctrl/internal/config"
)

type Daemon struct {
	address    string
	layout     *api.Layout
	logHandler slog.Handler
}

type handlerFunc[T any] func(*Daemon, *http.Request) (T, error)

func writeJSON[T any](w http.ResponseWriter, status int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return fmt.Errorf("failed encoding response: %w", err)
	}

	return nil
}

func handle[T any](d *Daemon, handler handlerFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var writeErr error

		defer func() {
			if writeErr != nil {
				d.Logger().Error(writeErr.Error())
			}
		}()

		result, err := handler(d, r)
		if err != nil {
			if apiErr, ok := errors.AsType[*api.APIError](err); ok {
				writeErr = writeJSON(w, apiErr.Code, map[string]string{
					"error": apiErr.Message,
				})
				if writeErr != nil {
					d.Logger().Error(writeErr.Error())
				}

				return
			}

			// Unexpected errors become 500.
			writeErr = writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			if writeErr != nil {
				d.Logger().Error(writeErr.Error())
			}

			return
		}

		writeErr = writeJSON(w, http.StatusOK, result)
	}
}

func (d *Daemon) endpoints() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/route/find", handle(d, getRouteFind))

	return mux
}

func (d *Daemon) Logger() *slog.Logger {
	return slog.New(d.logHandler)
}

func (d *Daemon) Start() error {
	server := &http.Server{
		Handler: d.endpoints(),
	}

	var listener net.Listener
	var err error

	if d.address == "" {
		socket := api.DaemonSocket

		// Cleanup socket in case it exists.
		_ = os.Remove(socket)

		listener, err = net.Listen("unix", socket)
	} else {
		listener, err = net.Listen("tcp", d.address)
	}

	if err != nil {
		return fmt.Errorf("failed starting listener: %w", err)
	}

	err = server.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed running API server: %w", err)
	}

	return nil
}

func NewDaemon(address string, configPath string) (*Daemon, error) {
	layout, err := config.DecodeLayout(configPath)
	if err != nil {
		return nil, err
	}

	logHandler := slog.NewTextHandler(os.Stdout, nil)

	return &Daemon{
		address:    address,
		layout:     layout,
		logHandler: logHandler,
	}, nil
}
