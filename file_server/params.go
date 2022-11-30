package main

import (
	"mime/multipart"
	"strconv"
)

func checkUUID(uuid string) bool {
	return uuid != "" && (len(uuid) == 32)
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

func checkAlltotal(total int64) bool {
	return total < MaxAllTotalSize
}

func checkMultiPartFileHeader(header *multipart.FileHeader) bool {
	return header.Size > 0 && header.Size < MaxSegmentSize
}
