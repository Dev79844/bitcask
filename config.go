package bitcask

import "time"

const (
	defaultDir 							= 	"./tmp"
	defaultMaxActiveFileSize 			= 	int64(1<<32) // 4GB
	defaultCompactInterval 				= 	time.Hour * 6
	defaultSyncInterval 				=	time.Minute * 1
	defaultActiveFileSizeCheckInterval 	= 	time.Minute * 1
)

type Options struct{
	dir						string 			//Path for storing active files
	alwaysFSync				bool 			//Always flush to disk after a write
	compactInterval 		time.Duration 	//time interval for compaction
	syncInterval    		*time.Duration 	// time interval for fsync
	maxActiveFileSize		int64 			//max active file size
	checkFileSizeInterval	time.Duration	// checks file size periodically
}

type Config func(*Options) error

func DefaultOptions() *Options {
	return &Options{
		dir: defaultDir,
		alwaysFSync: false,
		compactInterval: defaultCompactInterval,
		maxActiveFileSize: defaultMaxActiveFileSize,
		checkFileSizeInterval: defaultActiveFileSizeCheckInterval,
	}
}

func WithDir(dir string) Config {
	return func(o *Options) error {
		o.dir = dir
		return nil
	}
}

func WithAlwaysFSync() Config {
	return func(o *Options) error {
		o.alwaysFSync = true
		return nil
	}
}

func WithCompactInterval(interval time.Duration) Config {
	return func(o *Options) error {
		o.compactInterval = interval
		return nil
	}
}

func WithSyncInterval(interval time.Duration) Config {
	return func(o *Options) error {
		o.alwaysFSync = false
		o.syncInterval = &interval
		return nil
	}
}

func WithMaxActiveFileSize(fileSize int64) Config {
	return func(o *Options) error {
		o.maxActiveFileSize = fileSize
		return nil
	}
}