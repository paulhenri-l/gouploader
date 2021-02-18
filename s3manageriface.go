//go:generate mockgen -package mocks -destination ./mocks/s3manageriface.go github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface UploaderAPI

package gouploader
