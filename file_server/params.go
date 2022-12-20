package main

import (
	"mime/multipart"
	"strconv"

	"github.com/cwloo/uploader/file_server/config"
)

func checkUUID(uuid string) bool {
	return uuid != "" && (len(uuid) == 36) &&
		uuid[8] == '-' && uuid[13] == '-' &&
		uuid[18] == '-' && uuid[23] == '-'
}

func checkMD5(md5 string) bool {
	return md5 != "" && (len(md5) == 32)
}

func checkSingle(total string) bool {
	if total == "" {
		return false
	}
	size, _ := strconv.ParseInt(total, 10, 0)
	if size <= 0 || size >= config.Config.Upload.MaxSingleSize {
		return false
	}
	return true
}

func checkOffset(offset, total string) bool {
	if offset == "" {
		return false
	}
	now, _ := strconv.ParseInt(offset, 10, 0)
	size, _ := strconv.ParseInt(total, 10, 0)
	if now < 0 || now >= size {
		return false
	}
	return true
}

func checkTotal(total int64) bool {
	return total < config.Config.Upload.MaxTotalSize
}

func checkMultiPartSize(header *multipart.FileHeader) bool {
	return header.Size > 0
}

func checkMultiPartSizeLimit(header *multipart.FileHeader) bool {
	return header.Size < config.Config.Upload.MaxSegmentSize
}
