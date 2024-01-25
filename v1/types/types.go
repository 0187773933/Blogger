package types

type ConfigFile struct {
	ServerBaseUrl string `yaml:"server_base_url"`
	ServerLiveUrl string `yaml:"server_live_url"`
	ServerPort string `yaml:"server_port"`
	ServerAPIKey string `yaml:"server_api_key"`
	ServerCookieName string `yaml:"server_cookie_name"`
	ServerCookieSecret string `yaml:"server_cookie_secret"`
	ServerCookieAdminSecretMessage string `yaml:"server_cookie_admin_secret_message"`
	ServerCookieSecretMessage string `yaml:"server_cookie_secret_message"`
	AdminUsername string `yaml:"admin_username"`
	AdminPassword string `yaml:"admin_password"`
	TimeZone string `yaml:"time_zone"`
	BoltDBPath string `yaml:"bolt_db_path"`
	BoltDBEncryptionKey string `yaml:"bolt_db_encryption_key"`
}

// type PostTypes  struct {
// 	Text string `yaml:"text"`
// 	HTML string `yaml:"html"`
// 	Bytes []byte `yaml:"bytes"`
// }

type Post struct {
	Date string `yaml:"date"`
	Type string `yaml:"type"`
	HTML string `yaml:"html"`
	Text string `yaml:"text"`
	MD []string `yaml:"mark_down"`
	Bytes []byte `yaml:"bytes"`
}