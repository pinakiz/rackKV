package rackkv

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DataDir       string
	HintDir       string
	MergeInterval int
	MaxFileSizeMB int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}

	cfg.DataDir = os.Getenv("DATA_DIR")
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}

	cfg.HintDir = os.Getenv("HINT_DIR")
	if cfg.DataDir == "" {
		cfg.DataDir = "./data"
	}

	mergeIntStr := os.Getenv("MERGE_INTERVAL")
	if mergeIntStr == "" {
		cfg.MergeInterval = 120 // default 120 sec
	} else {
		val, err := strconv.Atoi(mergeIntStr)
		if err != nil {
			return nil, fmt.Errorf("invalid MERGE_INTERVAL: %w", err)
		}
		cfg.MergeInterval = val
	}

	maxFileStr := os.Getenv("MAX_FILE_SIZE_MB")
	if maxFileStr == "" {
		cfg.MaxFileSizeMB = 10
	} else {
		val, err := strconv.Atoi(maxFileStr)
		if err != nil {
			return nil, fmt.Errorf("invalid MAX_FILE_SIZE_MB: %w", err)
		}
		cfg.MaxFileSizeMB = val
	}

	return cfg, nil
}
