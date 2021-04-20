package entities

import "time"

type UploadResult struct {
	Filepath string
	Size     int64
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Error    error
}

func (ur *UploadResult) GetDuration() time.Duration {
	return ur.End.Sub(ur.Start)
}

func (ur *UploadResult) GetSpeed() float64 {
	size := float64(ur.Size / 1024 / 1024)

	return size / ur.GetDuration().Seconds()
}

func (ur *UploadResult) GetError() error {
	return ur.Error
}
