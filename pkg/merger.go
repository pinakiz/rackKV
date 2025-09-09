package pkg

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	rackkv "rackKV"
	"sort"
	"strings"
	"time"
)

func Merger(in_mem_map map[string]KeyDirEntry, handler *RackHandle) error {
	fmt.Println("merging...")

	dataDir := "./data"

	// Clean up old temp files
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("error reading data directory: %w", err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".temp") {
			_ = os.Remove(filepath.Join(dataDir, f.Name()))
		}
	}

	// Collect .data files
	files, err = os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("error reading data directory after cleanup: %w", err)
	}
	var oldFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".data") {
			oldFiles = append(oldFiles, f.Name())
		}
	}
	sort.Slice(oldFiles, func(i, j int) bool {
		idi, _ := File_name_to_Id(oldFiles[i])
		idj, _ := File_name_to_Id(oldFiles[j])
		return idi < idj
	})

	// Find starting id
	activeId, err := GetActiveFile(dataDir)
	if err != nil {
		activeId = 0
	}
	nextId := activeId + 1

	// Open first temp file
	tempName := Id_to_file_name(nextId)
	tempPath := filepath.Join(dataDir, tempName)
	newFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("error creating temp file %s: %w", tempPath, err)
	}
	defer func() {
		if newFile != nil {
			_ = newFile.Close()
		}
	}()

	var written int64
	const headerSize = 20 // 4+8+4+4

	// write
	writeRecord := func(f *os.File, crc uint32, ts uint64, keySz, valSz uint32, keyBuf, valBuf []byte) error {
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, crc)
		binary.Write(&buf, binary.LittleEndian, ts)
		binary.Write(&buf, binary.LittleEndian, keySz)
		binary.Write(&buf, binary.LittleEndian, valSz)
		buf.Write(keyBuf)
		buf.Write(valBuf)

		handler.mu.Lock()
		_, err := f.Write(buf.Bytes())
		handler.mu.Unlock()
		if err != nil {
			return err
		}
		written += int64(buf.Len())
		return nil
	}

	// rotate to next .temp
	rollFile := func() error {
		if newFile != nil {
			_ = newFile.Sync()
			_ = newFile.Close()
			newFile = nil
		}
		// rename current temp -> data
		dataName := Id_to_file_name(nextId)
		if err := os.Rename(filepath.Join(dataDir, Id_to_file_name(nextId)), filepath.Join(dataDir, dataName)); err != nil {
			return fmt.Errorf("error renaming temp to data: %w", err)
		}
		// prepare next temp
		nextId++
		tempName = Id_to_file_name(nextId)
		tempPath = filepath.Join(dataDir, tempName)
		var err error
		newFile, err = os.OpenFile(tempPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}
		written = 0
		return nil
	}

	// Process old files
	for _, oldFile := range oldFiles {
		oldPath := filepath.Join(dataDir, oldFile)
		readFile, err := os.Open(oldPath)
		if err != nil {
			return fmt.Errorf("error opening %s: %w", oldPath, err)
		}

		for {
			header := make([]byte, headerSize)
			_, err := io.ReadFull(readFile, header)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			if err != nil {
				readFile.Close()
				return fmt.Errorf("error reading header from %s: %w", oldPath, err)
			}

			crc := binary.LittleEndian.Uint32(header[0:4])
			ts := binary.LittleEndian.Uint64(header[4:12])
			keySz := binary.LittleEndian.Uint32(header[12:16])
			valSz := binary.LittleEndian.Uint32(header[16:20])

			keyBuf := make([]byte, keySz)
			if _, err := io.ReadFull(readFile, keyBuf); err != nil {
				break
			}
			valBuf := make([]byte, valSz)
			if _, err := io.ReadFull(readFile, valBuf); err != nil {
				break
			}

			key := string(keyBuf)
			keyEntry, ok := in_mem_map[key]
			if !ok {
				continue
			}

			// calculate value position for this record
			curOffset, _ := readFile.Seek(0, io.SeekCurrent)
			recordStart := curOffset - int64(headerSize) - int64(keySz) - int64(valSz)
			valuePos := recordStart + int64(headerSize) + int64(keySz)

			if valuePos != keyEntry.ValuePos {
				continue
			}

			// CRC check
			checkBuf := append(header[4:20], keyBuf...)
			checkBuf = append(checkBuf, valBuf...)
			if crc32.ChecksumIEEE(checkBuf) != crc {
				fmt.Println("CRC mismatch for key:", key)
				continue
			}

			// write to new file
			if err := writeRecord(newFile, crc, ts, keySz, valSz, keyBuf, valBuf); err != nil {
				readFile.Close()
				return fmt.Errorf("error writing record: %w", err)
			}

			if written > 10*1024*1024 {
				if err := rollFile(); err != nil {
					readFile.Close()
					return err
				}
			}
		}
		readFile.Close()
		_ = os.Remove(oldPath)
	}

	if newFile != nil {
		_ = newFile.Sync()
		_ = newFile.Close()
		finalData := Id_to_file_name(nextId)
		if err := os.Rename(filepath.Join(dataDir, Id_to_file_name(nextId)), filepath.Join(dataDir, finalData)); err != nil {
			fmt.Println("Warning: failed to rename last temp file", err)
		}
	}
	fmt.Println("Generating Hint Files");
	err = Generate_hintFiles();
	if err != nil{
		return fmt.Errorf("error while generating hint files after merging: %w",err);
	}
	fmt.Println("Merge complete.")
	return nil
}
func MergerListener(ctx context.Context, in_mem_map map[string]KeyDirEntry, handler *RackHandle) {
	cfg,err := rackkv.LoadConfig();
	if(err != nil){
		 fmt.Println("error while loading merging: %w",err);
	}
	ticker := time.NewTicker(time.Duration(cfg.MergeInterval)* time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := Merger(in_mem_map, handler); err != nil {
				fmt.Println("error during merger:", err)
			} else {
				fmt.Println("merger completed successfully")
			}
		case <-ctx.Done():
			fmt.Println("MergerListener stopped")
			return
		}
	}
}
