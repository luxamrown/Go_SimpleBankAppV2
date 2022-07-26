package config

import (
	"fmt"

	"github.com/spf13/viper"
	"mohamadelabror.me/simplebankappv2/manager"
)

type ApiConfig struct {
	Url string
}

type DbConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Manager struct {
	InfraManager   manager.Infra
	RepoManager    manager.RepoManager
	UseCaseManager manager.UseCaseManager
}

type Config struct {
	Manager
	DbConfig
	ApiConfig
}

func (c Config) ReadConfigFile(path string, name string) Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigName(name)
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		return Config{}
	}

	c.ApiConfig = ApiConfig{Url: v.GetString("API_URL_DEPLOY")}
	c.DbConfig = DbConfig{
		Host:     v.GetString("DB_HOST"),
		Port:     v.GetString("DB_PORT"),
		Name:     v.GetString("DB_NAME"),
		User:     v.GetString("DB_USER"),
		Password: v.GetString("DB_PASSWORD"),
	}
	return c
}

func NewConfig(path, configFileName string) Config {
	cfg := Config{}
	cfg = cfg.ReadConfigFile(path, configFileName)

	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbConfig.User, cfg.DbConfig.Password, cfg.DbConfig.Host, cfg.DbConfig.Port, cfg.DbConfig.Name)

	cfg.InfraManager = manager.NewInfra(dataSourceName)
	cfg.RepoManager = manager.NewRepoManager(cfg.InfraManager)
	cfg.UseCaseManager = manager.NewUseCaseManager(cfg.RepoManager)
	return cfg
}
