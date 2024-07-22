package s3

type Config struct {
	URL             string `json:"url"`
	AccessKey       string `json:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey"`
	SSL             bool   `json:"ssl"`
	BucketName      string
}
