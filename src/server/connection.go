package server

import "sync/atomic"

type Connection struct {
	Active   uint64 `json:"active"`
	Executed uint64 `json:"executed"`
	Errors   uint64 `json:"errors"`
}

func (c *Connection) Inc() {
	atomic.AddUint64(&c.Active, 1)
	atomic.AddUint64(&c.Executed, 1)
}

func (c *Connection) Dec() {
	if c.Active > 0{
		atomic.StoreUint64(&c.Active, atomic.LoadUint64(&c.Active)-1)
	}
}

func (c *Connection) Error() {
	atomic.AddUint64(&c.Errors, 1)
}
