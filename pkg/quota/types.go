package quota

type QuotaReport struct {
	ID    uint32
	Used  uint64
	Limit uint64
}
