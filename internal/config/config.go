package config

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DATA_BASE_URI string `json:"data_base_uri"`
	Port          string `json:"port"`
}
type ConfigStorage struct {
	DB *pgxpool.Pool
}

func GetConfigFromENV() string {
	var cfgStore ConfigStorage
	if err := envconfig.Process("APP", &cfgStore); err != nil {
		return ""
	}
	return cfgStore.DB.Config().ConnConfig.Database

}
func NewConfigStorage() *ConfigStorage {
	str := GetConfigFromENV()
	pool, err := pgxpool.New(context.Background(), str)
	if err != nil {
		return nil
	}
	return &ConfigStorage{
		DB: pool,
	}
}
