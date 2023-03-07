package models

type Config struct {
	Address     string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DataBaseURI string `env:"DATABASE_URI"`
	SecretKey   string `env:"SECRET_KEY" envDefault:"secret-key"`
}
