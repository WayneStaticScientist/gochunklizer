package chunkedupload

import (
	"sync"

	"opechains.shop/chunklizer/v2/types"
)

var chunkChan = make(chan types.ChunkCache)
var chunkCache = make(map[string]types.ChunkCache)
var chunkCacheMutex sync.RWMutex

type ChunkUploader struct {
}

func InitChunkUploader() ChunkUploader {
	return ChunkUploader{}
}
