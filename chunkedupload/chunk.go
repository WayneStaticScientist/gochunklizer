package chunkedupload

import (
	"sync"

	"opechains.shop/chunklizer/v2/types"
	"opechains.shop/chunklizer/v2/websocket"
)

var chunkCacheMutex sync.RWMutex
var chunkChan = make(chan types.ChunkCache)
var chunkCache = make(map[string]types.ChunkCache)

type ChunkUploader struct {
	Socket *websocket.WebSocketManager
}

func InitChunkUploader(socket *websocket.WebSocketManager) ChunkUploader {
	return ChunkUploader{
		Socket: socket,
	}
}

// /verifyToken /user/verifyToken?t=
