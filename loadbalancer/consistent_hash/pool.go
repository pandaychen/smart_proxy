package consistent_hash

import "sync"

type ConsistentHashPool struct {
	sync.RWMutex
	vNodes map[uint32]*ConsistenthashBackendNodeWrapper
}
