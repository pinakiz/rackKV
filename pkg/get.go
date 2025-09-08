package pkg

import (
	"fmt"
	"io"
	"os"
)


func GET(handler *RackHandle, key string) (string, error) {
    keyEntry, ok := handler.KeyDir[key]
    if !ok {
        return "", fmt.Errorf("key not found: %s", key)
    }

    fileName := Id_to_file_name(keyEntry.FileId)

    f, err := os.Open("./data/" + fileName)
    if err != nil {
        return "", fmt.Errorf("error while opening data file: %w", err)
    }
    defer f.Close()

    // Seek to the value position
    if _, err := f.Seek(keyEntry.ValuePos, io.SeekStart); err != nil {
        return "", fmt.Errorf("seek failed: %w", err)
    }

    // Read the value
    valBuf := make([]byte, keyEntry.ValueSz)
    if _, err := io.ReadFull(f, valBuf); err != nil {
        return "", fmt.Errorf("error while reading value: %w", err)
    }

    return string(valBuf), nil
}
