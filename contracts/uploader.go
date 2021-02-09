//go:generate mockgen -package contracts -destination ../mocks/contracts/uploader.go . Uploader

package contracts

type Uploader interface {
	Upload(filepath string) UploadResult
}
