package lock

const (
	lockPath          = "/locks"
	readerLockPath    = "/readerLock"
	writerLockPath    = "/writerLock"
	ReaderNumRootPath = "/readers"
	globalLockPath = "/globalLock"
	globalReaderPath = "/global" /* i.e., /readers/global */
)
