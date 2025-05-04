package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Google   GoogleConfig   `mapstructure:"google"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port int    `mapstructure:"port"`
}

type GoogleConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func LoadConfig() (*Config, error) {
	viper.BindEnv("google.client_id", "GOOGLE_CLIENT_ID")
    viper.BindEnv("google.client_secret", "GOOGLE_CLIENT_SECRET")
    viper.BindEnv("google.redirect_url", "REDIRECT_URL")
    viper.BindEnv("postgres.user", "POSTGRES_USER")
    viper.BindEnv("postgres.password", "POSTGRES_PASSWORD")
    viper.BindEnv("postgres.dbname", "POSTGRES_DB")
    viper.BindEnv("postgres.host", "POSTGRES_HOST")
    viper.BindEnv("postgres.port", "POSTGRES_PORT") 
    viper.BindEnv("postgres.sslmode", "POSTGRES_SSLMODE") 
    viper.BindEnv("redis.password", "REDIS_PASSWORD")
    viper.BindEnv("redis.host", "REDIS_HOST")       
    viper.BindEnv("redis.port", "REDIS_PORT")       
    viper.BindEnv("redis.db", "REDIS_DB")        
    viper.BindEnv("app.env", "APP_ENV")
    viper.BindEnv("app.port", "APP_PORT")         


	viper.SetDefault("postgres.host", "postgres")
	viper.SetDefault("postgres.port", "5432")
	viper.SetDefault("postgres.sslmode", "disable")
	viper.SetDefault("redis.host", "redis")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("app.port", 8081)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		log.Println("Config file not found, relying on environment variables.")
	}

	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	if config.Google.ClientID == "" {
		log.Println("Warning: GOOGLE_CLIENT_ID is not set.")
	}
	if config.Google.ClientSecret == "" {
		log.Println("Warning: GOOGLE_CLIENT_SECRET is not set.")
	}
	log.Printf("DEBUG: Config - Postgres Host: [%s]", config.Postgres.Host)

	return &config, nil
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
