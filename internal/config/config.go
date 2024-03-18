package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Server          Server          `yaml:"server"`
	Database        DataBase        `yaml:"database"`
	Identity        Identity        `yaml:"identity"`
	FilmValidations FilmValidations `yaml:"filmValidations"`
}

type Server struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	Secret string `env:"SERVER_SECRET" env-required:"true"`
}

type DataBase struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	SSLMode  string `yaml:"sslmode"`
}

type Identity struct {
	MinLoginLen    int `yaml:"minLoginLen"`
	MaxLoginLen    int `yaml:"maxLoginLen"`
	MinPasswordLen int `yaml:"minPasswordLen"`
	MaxPasswordLen int `yaml:"maxPasswordLen"`
}

type FilmValidations struct {
	MinNameLen        int `yaml:"minNameLen"`
	MaxNameLen        int `yaml:"maxNameLen"`
	MinDescriptionLen int `yaml:"minDescriptionLen"`
	MaxDescriptionLen int `yaml:"maxDescriptionLen"`
	MinRating         int `yaml:"minRating"`
	MaxRating         int `yaml:"maxRating"`
}

func New(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
