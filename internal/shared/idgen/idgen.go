package idgen

type IDGenerator interface {
	NextID() (uint64, error)
}
