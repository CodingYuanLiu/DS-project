package lock

const (
	LockPath          = "/locks"
	ReaderLockPath    = "/readerLock"
	WriterLockPath    = "/writerLock"
	ReaderNumRootPath = "/readers"
	GlobalLockPath    = "/globalLock"
	GlobalReaderPath  = "/global" /* i.e., /readers/global */
)
