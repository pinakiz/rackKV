package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"time"
)

func computeCRC(timestamp int64, key []byte, val []byte) (uint32, error) {
    buf := new(bytes.Buffer)

    // Add timestamp
    if err := binary.Write(buf, binary.LittleEndian, timestamp); err != nil {
        return 0, err
    }

    // Add key length
    if err := binary.Write(buf, binary.LittleEndian, int32(len(key))); err != nil {
        return 0, err
    }

    // Add value length
    if err := binary.Write(buf, binary.LittleEndian, int32(len(val))); err != nil {
        return 0, err
    }

    // Add key bytes
    if _, err := buf.Write(key); err != nil {
        return 0, err
    }

    // Add value bytes
    if _, err := buf.Write(val); err != nil {
        return 0, err
    }

    // Compute CRC
    crc := crc32.ChecksumIEEE(buf.Bytes())
    return crc, nil
}

func PUT(handler *RackHandle , key string , val string)(string,error){
	activeFile := (handler.ActiveFileId)
	activeFileFD := handler.ActiveFile
	fileInfo, err := activeFileFD.Stat();
	
	if(err != nil){
		return "", fmt.Errorf("error getting stats of active file: %w",err);
	}
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if(fileInfo.Size()>10485760){
		newId,err  := GetActiveFile("./data");
		if(err!=nil){
			return "",fmt.Errorf("error getting new active id: %w",err);
		}
		newFile, err := os.OpenFile("./data/" +(Id_to_file_name(newId)),os.O_CREATE|os.O_RDWR,0666);
		if err != nil{
			return "", fmt.Errorf("error accessing new active file: %w",err);
		}
		handler.ActiveFileId = newId;
		handler.ActiveFile = newFile
	}

	if activeFileFD == nil {
		return "", fmt.Errorf("active file is nil")
	}
	// fmt.Println("active file : ",activeFile)
	tmstmp := time.Now().Unix();
	crc , err:= computeCRC(tmstmp,[]byte(key),[]byte(val));
	if(err != nil){
		return "",fmt.Errorf("error while computing checksum: %w",err);
	}
	buf := new(bytes.Buffer)

    if err := binary.Write(buf, binary.LittleEndian, crc); err != nil { return "", err }
    if err := binary.Write(buf, binary.LittleEndian, tmstmp); err != nil { return "", err }
    if err := binary.Write(buf, binary.LittleEndian, int32(len(key))); err != nil { return "", err }
    if err := binary.Write(buf, binary.LittleEndian, int32(len(val))); err != nil { return "", err }
	buf.Write([]byte(key))
	buf.Write([]byte(val))

	headerSz := 4 + 8 + 4 + 4 // crc + timestamp + keySz + valueSz
    fileOffset, _ := activeFileFD.Seek(0, io.SeekEnd)
    valuePos := fileOffset + int64(headerSz) + int64(len(key))
	// fmt.Println("active file:",handler.ActiveFileId)
	entry := KeyDirEntry{
		FileId: int64(activeFile),
		ValueSz: int64(len(val)),
		ValuePos: valuePos,
		Tstamp: tmstmp,
	}
	handler.KeyDir[key] = entry;
	_, err = activeFileFD.Write(buf.Bytes())
	handler.WriteCount++
	if handler.WriteCount%100 == 0 || time.Since(handler.LastSync) > 100*time.Millisecond {
		if err := handler.ActiveFile.Sync(); err != nil {
			return "", fmt.Errorf("sync error: %w", err)
		}
		handler.LastSync = time.Now()
		handler.WriteCount = 0
	}

	if(err!=nil){
		return "" ,fmt.Errorf("can't put the entry");
	}else{
		return "Ok" , nil
	}
}
