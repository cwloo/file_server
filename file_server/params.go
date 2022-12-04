package main

import (
	"mime/multipart"
	"strconv"
)

func checkUUID(uuid string) bool {
	return uuid != "" && (len(uuid) == 36) &&
		uuid[8] == '-' && uuid[13] == '-' &&
		uuid[18] == '-' && uuid[23] == '-'
}

func checkMD5(md5 string) bool {
	return md5 != "" && (len(md5) == 32)
}

func checkTotal(total string) bool {
	if total == "" {
		return false
	}
	size, _ := strconv.ParseInt(total, 10, 0)
	if size <= 0 || size >= MaxTotalSize {
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

func checkAlltotal(total int64) bool {
	return total < MaxAllTotalSize
}

func checkMultiPartSizeLimit(header *multipart.FileHeader) bool {
	return header.Size < MaxSegmentSize
}

func checkMultiPartSize(header *multipart.FileHeader) bool {
	return header.Size > 0
}
