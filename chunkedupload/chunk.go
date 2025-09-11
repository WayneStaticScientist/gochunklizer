package chunkedupload

import (
	"sync"

	"opechains.shop/chunklizer/v2/types"
)

var chunkCacheMutex sync.RWMutex
var chunkChan = make(chan types.ChunkCache)
var chunkCache = make(map[string]types.ChunkCache)

type ChunkUploader struct {
}

func InitChunkUploader() ChunkUploader {
	return ChunkUploader{}
}

// /verifyToken /user/verifyToken?t=
