package contracts

import "time"

type UploadResult interface {
	GetFilepath() string
	GetSize() int64
	GetStart() time.Time
	GetEnd() time.Time
	GetDuration() time.Duration
	GetSpeed() float64
	GetError() error
}
