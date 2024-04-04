package config

type ExternalAPIConfig struct {
	UserApiUrl string
}

func getExternalAPIConfig() ExternalAPIConfig {
	return ExternalAPIConfig{
		UserApiUrl: getString("USER_API_URL"),
	}
}
