package pkg

import (
	"fmt"
	"io"
	"os"
)


func GET(handler *RackHandle, key string) (string, error) {
	keyEntry := handler.KeyDir[key];
	fileName := Id_to_file_name(int64(keyEntry.FileId));

	f , err := os.Open(fileName);
	if(err != nil){
		return "",fmt.Errorf("error while reading the data files: %w",err);
	}
	defer f.Close();

	valBuf := make([]byte , keyEntry.ValueSz);
	val,err := io.ReadFull(f , valBuf);
	
	if(err!=nil){
		return "",fmt.Errorf("erro while reading the buffer: %w",err);
	}
	return string(val),nil;

}