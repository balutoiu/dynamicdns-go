package mailinabox

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Domain   string `yaml:"domain"`
	APIUrl   string `yaml:"api_url"`
}
