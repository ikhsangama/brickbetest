package config

func BankBaseUrl() string {
	const defaultUrl = "https://65e7db9d53d564627a8f5a3a.mockapi.io"

	return getEnvWithDefault("BANK_BASE_URL", defaultUrl)
}
