package idgen

import (
	"errors"

	"github.com/algorithm9/flash-deal/internal/model"
)

type SnowflakeIDGen struct {
	sf *Snowflake
}

func (g *SnowflakeIDGen) NextID() (uint64, error) {
	return g.sf.NextID()
}

func NewSnowflakeIDGen(machine *model.Machine) (*SnowflakeIDGen, error) {
	sf := NewSnowflake(int64(machine.ID))
	if sf == nil {
		return nil, errors.New("failed to initialize Sonyflake")
	}
	return &SnowflakeIDGen{sf: sf}, nil
}
