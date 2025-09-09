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
	SyncEveryN 	int
	SyncInterval	int
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

	syncEveryN := os.Getenv("SYNC_EVERY_N_STEPS")
	if syncEveryN != "" {
		val, err := strconv.Atoi(syncEveryN)
		if err != nil {
			return nil, fmt.Errorf("invalid SYNC_EVERY_N_STEPS: %w", err)
		}
		cfg.SyncEveryN = val
	}else{
		cfg.SyncEveryN = 100
	}

	syncInterval := os.Getenv("SYNC_INTERVAL")
	if syncInterval != "" {
		val, err := strconv.Atoi(syncInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid SYNC_INTERVAL: %w", err)
		}
		cfg.SyncInterval = val
	}else {
		cfg.SyncInterval = 10; 
	}

	return cfg, nil
}
