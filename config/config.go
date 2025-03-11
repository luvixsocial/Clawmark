package config

type Config struct {
	Server   Server   `yaml:"server" validate:"required"`
	Database Database `yaml:"storage" validate:"required"`
}

type Server struct {
	Port string `yaml:"port" comment:"Server Port" validate:"required"`
	Env  string `yaml:"env" comment:"Server Environment" validate:"required"`
}

type Database struct {
	DatabaseURL string `yaml:"database_url" comment:"Database URL" validate:"required"`
	RedisURL    string `yaml:"redis_url" comment:"Redis URL" validate:"required"`
}
