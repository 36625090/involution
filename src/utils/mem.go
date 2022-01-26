package utils

import (
	"runtime"
	"time"
)

func MemStats() map[string]interface{}  {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	return map[string]interface{}{
		"alloc":          mem.Alloc,
		"total_alloc":    mem.TotalAlloc,
		"heap_alloc": mem.HeapAlloc,
		"lookups":        mem.Lookups,
		"pause_total_ns": mem.PauseTotalNs,
		"enable_gc":      mem.EnableGC,
		"num_gc":         mem.NumGC,
		"num_force_gc":   mem.NumForcedGC,
		"next_gc":        mem.NextGC,
		"last_gc_at_ms":        mem.LastGC / uint64(time.Millisecond),
	}
}
