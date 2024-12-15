package quack

import (
	"time"
)

type V interface {
	bool | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string | time.Time | time.Duration
}
