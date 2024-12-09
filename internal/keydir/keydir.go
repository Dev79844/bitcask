package keydir

type KeyDir map[string]Meta

type Meta struct {
	FileID		int
	RecordSz	int
	RecordPos	int
	Tstamp		int
}