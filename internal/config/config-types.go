package config

type Config struct {
	Redis *Redis
}

type Redis struct {
	Host     string
	Password string
}
