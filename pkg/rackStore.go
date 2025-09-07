package pkg

import(
	"os"
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

type RackHandle{
	ActiveFileId int64
	LOCKFILE 
	ActiveFile *os.File();
	Mode Mode
	
}

func Open(){

}