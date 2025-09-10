package chunkedupload

import (
	"opechains.shop/chunklizer/v2/database"
)

type ChunkUploader struct {
	db database.Database
}

func InitChunkUploader(db database.Database) ChunkUploader {
	return ChunkUploader{
		db: db,
	}
}
