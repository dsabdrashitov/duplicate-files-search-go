package csvdb

const (
	ErroneousRow = ",\"e\"e\x0A" // 6 character string that breakes parsing from any state
)

const (
	fileBufferSize    = 512
	serviceMultiplier = 3
)
