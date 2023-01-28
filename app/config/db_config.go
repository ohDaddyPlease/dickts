package config

type DB struct {
	Port         string `yaml:"port"`
	Address      string `yaml:"address"`
	UserName     string `yaml:"user_name"`
	Password     string `yaml:"password"`
	DatabaseName string `yaml:"database_name"`
}
