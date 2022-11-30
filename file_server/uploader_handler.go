package main

import "time"

// <summary>
// Uploader
// <summary>
type Uploader interface {
	Get() time.Time
	Upload(req *Req)
	Clear()
	Close()
	NotifyClose()
}
