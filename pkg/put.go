package pkg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"time"
)

// need more work
// activefile is not correct
func PUT(handler *RackHandle , key string , val string)(string,error){
	activeFile := (handler.ActiveFileId)
	activeFileFD := handler.ActiveFile
	if activeFileFD == nil {
    return "", fmt.Errorf("active file is nil")
}

	fmt.Println("active file : ",activeFile)

	crc := crc32.ChecksumIEEE([]byte(key + val))
	buf := new(bytes.Buffer)

    if err := binary.Write(buf, binary.LittleEndian, crc); err != nil { return "", err }
    if err := binary.Write(buf, binary.LittleEndian, time.Now().Unix()); err != nil { return "", err }
    if err := binary.Write(buf, binary.LittleEndian, int32(len(key))); err != nil { return "", err }
    if err := binary.Write(buf, binary.LittleEndian, int32(len(val))); err != nil { return "", err }
	buf.Write([]byte(key))
	buf.Write([]byte(val))

	headerSz := 4 + 8 + 4 + 4 // crc + timestamp + keySz + valueSz
	pos, _ := activeFileFD.Seek(0,io.SeekEnd);
	pos += int64(headerSz) + int64(len(key));

	entry := KeyDirEntry{
		FileId: int64(activeFile),
		ValueSz: int64(len(val)),
		ValuePos: pos,
		Tstamp: time.Now().Unix(),
	}
	handler.KeyDir[key] = entry;
	_, err := activeFileFD.Write(buf.Bytes())
	activeFileFD.Sync() // flush to disk

	if(err!=nil){
		return "" ,fmt.Errorf("can't put the entry");
	}else{
		return "Ok" , nil
	}
}