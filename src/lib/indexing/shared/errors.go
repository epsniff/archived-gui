package shared

import "fmt"

var (
	ErrInvalidIndex = fmt.Errorf("invalid index name")
	ErrUnknownIndex = fmt.Errorf("unknown index name")
	ErrIndexExists  = fmt.Errorf("index exits")

	ErrInvalidDocId = fmt.Errorf("invalid document id")
)
