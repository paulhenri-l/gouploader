//go:generate mockgen -package contracts -destination ../mocks/contracts/s3manageriface.go github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface UploaderAPI

package contracts
