package aqara

type AqaraConfig struct {
	ApiUrl      string `mapstructure:"api_url"`
	AppId       string `mapstructure:"app_id"`
	AccessToken string `mapstructure:"access_token"`
	KeyId       string `mapstructure:"key_id"`
	AppKey      string `mapstructure:"app_key"`
}

func NewConfig(ApiUrl, AppId, AccessToken, KeyId, AppKey string) AqaraConfig {
	return AqaraConfig{
		ApiUrl:      ApiUrl,
		AppId:       AppId,
		AccessToken: AccessToken,
		KeyId:       KeyId,
		AppKey:      AppKey,
	}
}
