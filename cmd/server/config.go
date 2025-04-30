package main

import (
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel   string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
	DB         Database   `yaml:"db"`
	GRPCServer GrpcServer `yaml:"grpc_server" env-prefix:"GRPC_"`
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

func MustLoad(cfgFilePath string) Config {
	var cfg Config

	err := cleanenv.ReadConfig(cfgFilePath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
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
