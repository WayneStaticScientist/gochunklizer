package chunkedupload

import "opechains.shop/chunklizer/v2/types"

var chunkCache = make(map[string]types.ChunkCache)

type ChunkUploader struct {
}

func InitChunkUploader() ChunkUploader {
	return ChunkUploader{}
}
