package pkg

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func ReadLogs(path string) error {
	dataDir := "./data"
	hintDir := "./hint"
	dataPath := filepath.Join(dataDir, path)
	tempHint := filepath.Join(hintDir, "temp.hint")

	readFile, err := os.Open(dataPath)
	if err != nil {
		return fmt.Errorf("cant open data file: %w", err)
	}
	defer readFile.Close()

	writeFile, err := os.OpenFile(tempHint, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("cant open hint file: %w", err)
	}
	defer writeFile.Close()

	reader := bufio.NewReader(readFile)

	for {
		// --- Read header ---
		header := make([]byte, 20) // CRC(4) + Timestamp(8) + KeySize(4) + ValueSize(4)
		_, err := io.ReadFull(reader, header)
		if err == io.EOF {
			break // finished reading file
		}
		if err != nil {
			return fmt.Errorf("error while reading header: %w", err)
		}

		crc := binary.LittleEndian.Uint32(header[0:4])
		keysz := binary.LittleEndian.Uint32(header[12:16])
		valsz := binary.LittleEndian.Uint32(header[16:20])

		// --- Read key + value ---
		keyBuf := make([]byte, keysz)
		if _, err := io.ReadFull(reader, keyBuf); err != nil {
			return fmt.Errorf("error while reading key: %w", err)
		}
		valBuf := make([]byte, valsz)
		if _, err := io.ReadFull(reader, valBuf); err != nil {
			return fmt.Errorf("error while reading value: %w", err)
		}

		// --- CRC check ---
		checkBuf := append(header[4:20], keyBuf...)
		checkBuf = append(checkBuf, valBuf...)
		if crc32.ChecksumIEEE(checkBuf) != crc {
			fmt.Printf("stored=%x computed=%x header=%x key=%x val=%x\n",
				crc, crc32.ChecksumIEEE(checkBuf), header, keyBuf, valBuf)

			return fmt.Errorf("CRC mismatch for key: %s", string(keyBuf))
		}

		// --- Calculate value position (relative to start of file) ---
		offset, _ := readFile.Seek(0, io.SeekCurrent)
		valPos := offset - int64(len(valBuf))

		// --- Build hint entry ---
		hintFileEntry := []byte{}
		hintFileEntry = append(hintFileEntry, header[4:]...)     // timestamp + keysize + valsize
		hintFileEntry = append(hintFileEntry, keyBuf...)         // key
		posBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(posBytes, uint64(valPos)) // val position
		hintFileEntry = append(hintFileEntry, posBytes...)

		if _, err := writeFile.Write(hintFileEntry); err != nil {
			return fmt.Errorf("error writing hint entry: %w", err)
		}
	}

	// --- Rename after processing all entries ---
	id, _ :=  File_name_to_Id(path)
	newHintName :=  Id_to_hint_name(id)
	if err := os.Rename(tempHint, filepath.Join(hintDir, newHintName)); err != nil {
		return fmt.Errorf("error renaming hint file: %w", err)
	}
	return nil
}

func Generate_hintFiles(handler *RackHandle) error {
	dataDir := "./data"
	hintDir := "./hint"
	activeFile := handler.ActiveFileId;

	dataFiles, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("error reading data dir: %w", err)
	}
	hintFiles, err := os.ReadDir(hintDir)
	if err != nil {
		return fmt.Errorf("error reading hint dir: %w", err)
	}

	dataFilesMap := make(map[int64]int64)

	for _, file := range dataFiles {
		fileID, err :=  File_name_to_Id(file.Name())
		if err != nil {
			return fmt.Errorf("error parsing data filename: %w", err)
		}
		dataFilesMap[fileID]++
	}

	for _, file := range hintFiles {
		fileID, err :=  File_name_to_Id(file.Name())
		if(fileID==activeFile){
			continue;
		}
		if err != nil {
			return fmt.Errorf("error parsing hint filename: %w", err)
		}
		if dataFilesMap[fileID]--; dataFilesMap[fileID] == 0 {
			delete(dataFilesMap, fileID)
		}
	}

	dataFilesList := make([]int64, 0, len(dataFilesMap))
	for fileID := range dataFilesMap {
		dataFilesList = append(dataFilesList, fileID)
	}

	sort.Slice(dataFilesList, func(i, j int) bool {
		return dataFilesList[i] > dataFilesList[j]
	})

	for _, fileID := range dataFilesList {
		dataFileName :=  Id_to_file_name(fileID)
		if err := ReadLogs(dataFileName); err != nil {
			return err
		}
	}

	return nil
}
