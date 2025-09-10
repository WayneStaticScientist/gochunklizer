package types

type ChunkCache struct {
	CurrentIndex int64
	TotalChunks  int64
	Step         int64
	ChunkPath    string
}
