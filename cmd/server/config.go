package main

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel   string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	DB         Database   `yaml:"db"`
	GRPCServer GrpcServer `yaml:"grpc_server" env-prefix:"GRPC_"`
	RabbitMQ   RabbitMQ   `yaml:"rabbitmq" env-prefix:"RABBITMQ_"`
}

type Database struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Name     string `yaml:"name" env:"DB_NAME" env-default:"postgres"`
	Username string `yaml:"username" env:"DB_USERNAME" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:""`
}

type GrpcServer struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"8082"`
}

type RabbitMQ struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"PORT" env-default:"5672"`
	Username string `yaml:"username" env:"USERNAME" env-default:"guest"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"guest"`
}

func MustLoad() Config {
	var cfg Config

	_ = godotenv.Load(".env")

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to read environment variables: %v", err)
	}

	return cfg
}

func (c *Config) MakeDBConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.DB.Username,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
		c.DB.Name,
	)
}

func (c *Config) MakeGRPCAddr() string {
	return fmt.Sprintf("%s:%d", c.GRPCServer.Host, c.GRPCServer.Port)
}
