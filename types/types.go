package types

type ChunkCache struct {
	CurrentIndex int64
	TotalChunks  int64
	Step         int64
	ChunkPath    string
	FileName     string
	FileType     string
	LastAccess   int64
	Token        string
	ObjectId     string
}
type UploadMetaDataFile struct {
	FileProvider string `json:"fileProvider"`
}
