package pkg

import (
	"fmt"
	"os"
)

func GetActiveFile(dir string) (int64, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return 0, fmt.Errorf("error while getting active file: %w", err)
	}

	var maxName int64

	for _, entries := range files {
		fileName, err := File_name_to_Id(entries.Name())
		if err != nil {
			return 0, fmt.Errorf("error while getting active file: %w", err)
		}
		if maxName < fileName {
			maxName = fileName
		}
	}
	return (maxName + 1), nil
}
