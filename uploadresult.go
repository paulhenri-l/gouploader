package gouploader

import "time"

type UploadResult struct {
	Filepath string
	Size     int64
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Error    error
}

func (ur *UploadResult) GetFilepath() string {
	return ur.Filepath
}

func (ur *UploadResult) GetSize() int64 {
	return ur.Size
}

func (ur *UploadResult) GetStart() time.Time {
	return ur.Start
}

func (ur *UploadResult) GetEnd() time.Time {
	return ur.End
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
