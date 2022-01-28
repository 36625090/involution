package server

import "sync/atomic"

type Connection struct {
	Active   int64 `json:"active"`
	Executed uint64 `json:"executed"`
	Errors   uint64 `json:"errors"`
}

func (c *Connection) Inc() {
	atomic.AddInt64(&c.Active, 1)
	atomic.AddUint64(&c.Executed, 1)
}

func (c *Connection) Dec() {
	if atomic.LoadInt64(&c.Active) > 0{
		atomic.AddInt64(&c.Active, -1)
	}
}

func (c *Connection) Error() {
	atomic.AddUint64(&c.Errors, 1)
}
