package shunt

import (
	"sync/atomic"
)

type Conns struct {
	connsindex map[uint32]*Conn
	connsuid   map[uint64]uint32
	index      uint32
}

var oneConns *Conns

func NewConns() *Conns {
	if oneConns == nil {
		oneConns = &Conns{
			connsindex: make(map[uint32]*Conn),
			connsuid:   make(map[uint64]uint32),
			index:      0,
		}
	}

	return oneConns
}

func (cs *Conns) Add(conn *Conn) {
	conn.index = atomic.AddUint32(&cs.index, 1)
	cs.connsindex[conn.index] = conn
	cs.connsuid[conn.IMEI] = conn.index
}

func (cs *Conns) GetConn(uid uint64) *Conn {
	return cs.connsindex[cs.connsuid[uid]]
}

func (cs *Conns) Remove(uid uint64) {
	index := cs.connsuid[uid]
	delete(cs.connsindex, index)
	delete(cs.connsuid, uid)
}

func (cs *Conns) Check(uid uint64) bool {
	_, ok := cs.connsindex[cs.connsuid[uid]]
	return ok
}

func (cs *Conns) GetCount() int {
	return len(cs.connsindex)
}
