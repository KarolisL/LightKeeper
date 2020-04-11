package daemon

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io"
)

type Filter struct {
	Program string
}

type Input struct {
	Name string
	Path string
	Type string
}

type Mapping struct {
	From    string
	To      string
	Filters []map[string]string
}

type Output struct {
	Type string
	Sink map[string]string
}

type Config struct {
	Inputs   map[string]Input
	Mappings []Mapping
	Outputs  map[string]Output
}

func NewConfigFromReader(r io.Reader) (*Config, error) {
	cfg := &Config{}
	err := toml.NewDecoder(r).Decode(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return cfg, nil
}
