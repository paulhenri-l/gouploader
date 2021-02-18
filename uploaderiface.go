//go:generate mockgen -package mocks -destination ./mocks/uploader_mock.go . Uploader

package gouploader

type Uploader interface {
	Upload(file string) *UploadResult
}
