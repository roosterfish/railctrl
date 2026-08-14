package config

import (
	"fmt"
	"os"
	"path/filepath"

	"githu.com/roosterfish/railctrl/api"
	"go.yaml.in/yaml/v4"
)

// DecodeLayout decodes the given layout config and returns the decoded contents.
// If an empty path is provided, we check the following in order:
// * $LAYOUT_CONFIG environment variable
// * ~/.railctrl/layout.yaml
func DecodeLayout(path string) (*api.Layout, error) {
	if path == "" {
		layoutEnv := os.Getenv("LAYOUT_CONFIG")
		if layoutEnv != "" {
			path = layoutEnv
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed getting user's home dir: %w", err)
			}

			path = filepath.Join(home, ".railctrl", "layout.yaml")
		}
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed reading from layout path %q: %w", path, err)
	}

	layout := &api.Layout{}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(layout)
	if err != nil {
		return nil, fmt.Errorf("failed decoding layout: %w", err)
	}

	return layout, nil
}

// EncodeLayout encodes the given layout and returns the encoded contents.
func EncodeLayout(layout *api.Layout) ([]byte, error) {
	encodedLayout, err := yaml.Marshal(layout)
	if err != nil {
		return nil, fmt.Errorf("failed encoding layout: %w", err)
	}

	return encodedLayout, nil
}
