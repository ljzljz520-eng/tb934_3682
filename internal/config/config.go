package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidConfig = errors.New("invalid configuration")

type Config struct {
	Address       string
	DatabasePath  string
	DefaultGuide  string
	ReaderLimit   int
	VisitorHeader string
}

func Default() Config {
	return Config{Address: ":8080", DatabasePath: "wedding-guide.db", DefaultGuide: "demo-guide", ReaderLimit: 4, VisitorHeader: "X-Visitor-Key"}
}

func FromArgs(args []string) (Config, error) {
	cfg := Default()
	for _, arg := range args {
		key, value, ok := strings.Cut(arg, "=")
		if !ok {
			return Config{}, fmt.Errorf("%w: %s", ErrInvalidConfig, arg)
		}
		switch strings.TrimLeft(key, "-") {
		case "address":
			cfg.Address = value
		case "db":
			cfg.DatabasePath = value
		case "guide":
			cfg.DefaultGuide = value
		case "reader-limit":
			parsed, err := strconv.Atoi(value)
			if err != nil {
				return Config{}, err
			}
			cfg.ReaderLimit = parsed
		case "visitor-header":
			cfg.VisitorHeader = value
		default:
			return Config{}, fmt.Errorf("%w: %s", ErrInvalidConfig, key)
		}
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Address) == "" || strings.TrimSpace(c.DatabasePath) == "" {
		return ErrInvalidConfig
	}
	if strings.TrimSpace(c.DefaultGuide) == "" || strings.TrimSpace(c.VisitorHeader) == "" {
		return ErrInvalidConfig
	}
	if c.ReaderLimit < 1 {
		return ErrInvalidConfig
	}
	return nil
}
