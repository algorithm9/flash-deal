package idgen

import (
	"errors"
	"sync"
	"time"
)

type Snowflake struct {
	mutex     sync.Mutex
	lastTs    int64
	sequence  int64
	machineID int64
}

// 常量配置
const (
	epoch         = int64(1735660800000000000) // 起始时间戳（单位：纳秒）
	machineIDBits = 10
	sequenceBits  = 12

	maxMachineID = -1 ^ (-1 << machineIDBits)
	maxSequence  = -1 ^ (-1 << sequenceBits)

	timeShift    = sequenceBits + machineIDBits
	machineShift = sequenceBits
)

// NewSnowflake 创建一个新 Snowflake 实例
func NewSnowflake(machineID int64) *Snowflake {
	if machineID < 0 || machineID > maxMachineID {
		panic("invalid machine id")
	}
	return &Snowflake{machineID: machineID}
}

// NextID 生成唯一 ID（时间回拨立即报错）
func (s *Snowflake) NextID() (uint64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixNano()
	if now < s.lastTs {
		// 时间回拨，立即报错
		return 0, errors.New("clock moved backwards: system time is behind last ID timestamp")
	}

	if now == s.lastTs {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 序列号用尽，不等待，直接报错
			return 0, errors.New("sequence overflow: too many IDs generated in the same millisecond")
		}
	} else {
		s.sequence = 0
	}

	s.lastTs = now

	id := ((now - epoch) >> 20 << timeShift) | // 纳秒转毫秒再左移（对齐位宽）
		(s.machineID << machineShift) |
		s.sequence

	return uint64(id), nil
}
