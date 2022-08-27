package config

type Config struct {
	Redis        *Redis
	Host         *Host
	DockerEngine *DockerEngine
}

type Redis struct {
	Host     string
	Password string
}

type Host struct {
	Port string
}

type DockerEngine struct {
	Imagename string
}
