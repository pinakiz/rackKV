package pkg

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// hint file record structure : [Timestamp][KeySize][ValueSize][Key][ValuePos]


func ReadHintFiles(path string , keydir map[string]KeyDirEntry) error {
	hintDir := "./hint/"+path
	readFile, err := os.Open(hintDir)
	if err != nil {
		return fmt.Errorf("cant open hint file: %w", err)
	}
	reader := bufio.NewReader(readFile)

	for {
		header := make([]byte, 24) // Timestamp(8) + KeySize(4) + ValueSize(4) + valuePos(8) + key(4);
		_, err := io.ReadFull(reader, header)
		if err == io.EOF {
			break // finished reading file
		}
		if err != nil {
			return fmt.Errorf("error while reading header: %w", err)
		}

		timestamp := binary.LittleEndian.Uint64(header[0:8])
		keysz := binary.LittleEndian.Uint32(header[8:12])
		valsz := binary.LittleEndian.Uint32(header[12:16])
		valPos := binary.LittleEndian.Uint64(header[16:24])
		keyBuf := make([]byte, keysz)
		if _, err := io.ReadFull(reader, keyBuf); err != nil {
			return fmt.Errorf("error while reading key: %w", err)
		}

		fileid , err := hint_name_to_Id(path);
		if(err != nil){
			return fmt.Errorf("error while converting hintId to hint_dir_name: %w" , err);
		} 
		temp := KeyDirEntry{
			FileId: fileid,
			ValueSz: int64(valsz),
			ValuePos: int64(valPos),
			Tstamp: int64(timestamp),
		}
		keydir[string(keyBuf)] = temp;

	}
	return nil
}



func GenerateKeyDir(handler *RackHandle)( error){
	hintDir := "./hint";
	hint_files , err := os.ReadDir(hintDir);
	if err != nil{
		 return fmt.Errorf("error while reading hint files: %w",err);
	}
	KeyDirEntry := make(map[string]KeyDirEntry);
	for _, files := range hint_files{

		if err := ReadHintFiles(files.Name() , KeyDirEntry);err != nil{
			return fmt.Errorf("error while reading hint files: %w",err);
		}
	}
	handler.KeyDir = KeyDirEntry;
	return nil;
	
}