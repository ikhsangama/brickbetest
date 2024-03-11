package config

func GetSqsUrl() string {
	const defaultUrl = "http://sqs.ap-southeast-1.localhost.localstack.cloud:4566/000000000000/"

	return getEnvWithDefault("SQS_URL", defaultUrl)
}
