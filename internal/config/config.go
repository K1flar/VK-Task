package config

type Config struct {
	DataBase
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}
