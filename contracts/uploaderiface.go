//go:generate mockgen -package mocks -destination ../mocks/uploader_mock.go . Uploader

package contracts

import "github.com/paulhenri-l/gouploader/entities"

type Uploader interface {
	Upload(file string) *entities.UploadResult
}
