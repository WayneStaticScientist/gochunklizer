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
	IsUploaded   bool
	IsFailed     bool
	IsPending    bool
	Progress     float64
	Message      string
	Tries        int64
}
type UploadMetaDataFile struct {
	FileProvider string `json:"fileProvider"`
}
type SocketChunckMessage struct {
	Data      any     `json:"data"`
	HasError  bool    `json:"isError"`
	Message   string  `json:"message"`
	Progress  float64 `json:"progress"`
	IsSuccess bool    `json:"isSuccess"`
}
