package config

type CertConfig struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
	Ca   string `yaml:"ca"`
}
