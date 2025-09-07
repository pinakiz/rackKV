package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

type KeyDirEntry struct{
	FileId int64
	ValueSz int64
	ValuePos int64
	Tstamp int64

}

type Mode struct{
	IsUp bool
	ReadWrite bool
	SyncOnWrite bool
}

type RackHandle struct{
	ActiveFileId int64
	lockFile *os.File 
	ActiveFile *os.File;
	Mode Mode
	KeyDir map[string]KeyDirEntry
	
}

func Open(directory string , mode Mode)(*RackHandle , error){
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Current Working Directory:", dir)

	handler := &RackHandle{}
	mode_handler := Mode{
		ReadWrite: mode.ReadWrite,
		SyncOnWrite: mode.SyncOnWrite,
	}
	handler.Mode = mode_handler

	handler.KeyDir = make(map[string]KeyDirEntry)

	activeFileId , err := GetActiveFile("../data");
	if(err != nil){
		return handler, fmt.Errorf("failed to get a active file: %w", err);
	}	
	name_activeFile := Id_to_file_name(activeFileId);

	activeFile , err := os.OpenFile("../data/"+name_activeFile , os.O_CREATE|os.O_RDWR , 0666);
	if(err != nil){	
		return handler, fmt.Errorf("failed to get a active data file: %w", err);
	}
	handler.ActiveFile = activeFile;

	if err := os.MkdirAll(directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create/open directory: %w", err)
	}
	lockFilePath := filepath.Join(directory,"LOCK");
	lockFile , err:= os.OpenFile(lockFilePath,os.O_CREATE|os.O_RDWR,0666);
	if(err != nil){
		return handler,nil;
	}
	handler.lockFile = lockFile;

	if(mode.ReadWrite){
		if err := syscall.Flock(int(lockFile.Fd()) , syscall.LOCK_EX|syscall.LOCK_NB); err!=nil{
			lockFile.Close();
			return handler,fmt.Errorf("error while aquaring the lock: %w",err);
		}
	}else{
		if err := syscall.Flock(int(lockFile.Fd()) , syscall.LOCK_SH|syscall.LOCK_NB); err!=nil{
			lockFile.Close();
			return handler,fmt.Errorf("error while aquaring the lock: %w",err);
		}
	}


	handler.Mode.IsUp = true;
	return handler,nil;
}

func (handler *RackHandle) Close() error{
	if err := syscall.Flock(int(handler.ActiveFile.Fd()),syscall.LOCK_UN); err != nil{
		return fmt.Errorf("failed to release the lock: %w",err);
	}else {
		fmt.Println("Datastore handle closed and lock released.")
		return nil;
	}
}