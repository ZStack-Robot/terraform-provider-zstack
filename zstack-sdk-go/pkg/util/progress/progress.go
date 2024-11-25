// Copyright (c) ZStack.io, Inc.

package progress

import (
	"io"
	"time"
)

func NewProgress(totalSize int64, maxPercent int, reader io.Reader, callback func(progress float32)) io.Reader {
	body := &Progress{
		total:      totalSize,
		maxPercent: maxPercent,
		callback:   callback,
	}
	return io.TeeReader(reader, body)
}

type Progress struct {
	refreshSeconds int
	count          int64
	total          int64
	start          time.Time
	callback       func(progress float32)
	maxPercent     int
}

func (r *Progress) Write(p []byte) (int, error) {
	if r.start.IsZero() {
		r.start = time.Now()
	}
	n := len(p)
	r.count += int64(n)
	if r.callback != nil && r.total > 0 && time.Now().Sub(r.start) > time.Second*1 {
		r.callback(float32(float64(r.count) / float64(r.total) * float64(r.maxPercent)))
		r.start = time.Now()
	}
	return n, nil
}
