package configuration

type Configuration struct {
	ExportBucketName string `long:"export-bucket-name" description:"The name of the S3 bucket to export change archives to"`
}
